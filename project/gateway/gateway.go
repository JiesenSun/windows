package gateway

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"project/common"
	"project/common/nsq"
	"project/common/protocol"
	sync_ "project/common/sync"
	"project/common/syslog"
	"project/common/util"
)

const (
	// timeout
	CONNECT_NO_LOGIN_TIMEOUT   = 30
	CONNECT_NO_REQUEST_TIMEOUT = 180
	// 客户端状态
	CLIENT_STATE_INIT    = 0
	CLIENT_STATE_ONLINE  = 1
	CLIENT_STATE_OFFLINE = 2

	// client num control
	CLIENT_NUM_PER_GROUP = 50
	CLIENT_GROUP_NUM     = 30
	PROGRAM_SCALE        = 15000
	// 读写时间
	NET_READ_DATELINE     = time.Millisecond * 20
	NET_READ_MAX_DATELINE = time.Second * 1
	NET_WRITE_DATELINE    = time.Second * 3
)

var (
	g_IP   uint32 = 0
	g_Port uint32 = 0
)

type Service struct {
	Listener         *net.TCPListener
	ClientCount      int32
	ClientCh         chan *Client
	ClientGroupCount int32
	WaitGroup        *sync_.WaitGroup
	ClientMap        map[uint32]*Client
	RWMutex          *sync.RWMutex
}

func (this *Service) HandleMessage(message *nsq.Message) error {
	dp := protocol.GetDataPackage()
	if err := dp.Unpackage(message.Body); err != nil {
		syslog.Warn(err, string(message.Body))
		protocol.PutDataPackage(dp)
		return nil
	}
	this.RWMutex.RLock()
	client, ok := this.ClientMap[dp.Tail.SID]
	this.RWMutex.RUnlock()

	if !ok || client.Stat == CLIENT_STATE_OFFLINE || dp.Tail.SID != client.Sid {
		protocol.PutDataPackage(dp)
		return nil
	}
	if client.RespCh == nil {
		protocol.PutDataPackage(dp)
		return nil
	}
	dp.Data = client
	client.RespCh <- dp
	atomic.AddInt32(&client.MsgCount, 1)
	return nil
}

func (this *Service) Serve() {
	this.WaitGroup.AddOne()
	defer this.WaitGroup.Done()
	exitNotify := this.WaitGroup.ExitNotify()
	var grpCount int32 = 0
	var clientCount int32 = 0
	var sid uint32 = 10000
	for {
		select {
		case <-exitNotify:
			return
		default:
		}

		this.Listener.SetDeadline(time.Now().Add(NET_READ_DATELINE))
		conn, err := this.Listener.AcceptTCP()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			syslog.Warn("Unknow Net Error:", err)
			break
		}
		syslog.Debug(conn.RemoteAddr(), "connected, seesionId:", sid)
		// 建立客户端
		client := GetClient()
		client.Service = this
		client.Conn = conn
		client.Sid = sid
		client.Stat = CLIENT_STATE_INIT
		client.RespCh = nil
		client.MsgCount = 0
		client.LastTime = time.Now().Unix()
		this.RWMutex.Lock()
		this.ClientMap[sid] = client
		this.RWMutex.Unlock()
		this.ClientCh <- client
		clientCount = atomic.AddInt32(&this.ClientCount, 1)

		grpCount = atomic.LoadInt32(&this.ClientGroupCount)
		if clientCount >= ((grpCount - 10) * CLIENT_NUM_PER_GROUP) {
			NewClientGroup(this)
		}
		// 处理客户端请求
		sid++ // increase sid
	}
}

func (this *Service) Stop() {
	this.WaitGroup.Wait()
	this.Listener.Close()
}

func (this *Service) P2pService(sock *net.UDPConn) {
	this.WaitGroup.AddOne()
	defer this.WaitGroup.Done()
	exitNotify := this.WaitGroup.ExitNotify()
	var buf [2048]byte
	for {
		select {
		case <-exitNotify:
			return
		default:
		}
		sock.SetDeadline(time.Now().Add(NET_READ_MAX_DATELINE))
		n, addr, err := sock.ReadFromUDP(buf[:])
		if err != nil {
			syslog.Info(err)
			continue
		}
		if n < 4 {
			syslog.Info("udp data package too small")
			continue
		}
		sid := binary.BigEndian.Uint32(buf[:4])
		this.RWMutex.RLock()
		client, ok := this.ClientMap[sid]
		this.RWMutex.RUnlock()
		if ok == false {
			syslog.Debug("client not exist!!! sid:", sid)
			continue
		}
		client.IP = util.IPToInt(addr.String())
		client.Port = uint32(addr.Port)
		if _, err := sock.WriteToUDP([]byte(addr.String()), addr); err != nil {
			syslog.Info("udp send data failed!!!", err)
			continue
		}
	}
}
func StartServer() {
	iniConf := common.Config()
	ip := iniConf.DefaultString("tcp_ip", "10.0.2.15")
	port := iniConf.DefaultString("tcp_port", "8888")
	g_IP = util.IPToInt(ip)
	g_Port = uint32(util.Atoi(port))

	listener, err := common.TCPServer(ip + ":" + port)
	if err != nil {
		panic(err)
	}
	/*
		udp_addr1 := iniConf.DefaultString("udp_address1", "127.0.0.1:9100")
		udp_addr2 := iniConf.DefaultString("udp_address2", "127.0.0.1:9101")

		udpSocket1, err := common.UDPServer(udp_addr1)
		if err != nil {
			panic(err)
		}
		udpSocket2, err := common.UDPServer(udp_addr2)
		if err != nil {
			panic(err)
		}
	*/
	service := &Service{
		Listener:  listener,
		ClientCh:  make(chan *Client, 1000),
		WaitGroup: sync_.NewWaitGroup(),
		ClientMap: make(map[uint32]*Client),
		RWMutex:   &sync.RWMutex{},
	}
	for i := 0; i < CLIENT_GROUP_NUM; i++ {
		NewClientGroup(service)
	}
	go service.Serve()
	//go service.P2pService(udpSocket1)
	//go service.P2pService(udpSocket2)

	topic := fmt.Sprintf("conn_%v_%v", ip, port)
	if _, err := common.NsqConsumerGO(topic, "conn-channel", 3, service); err != nil {
		panic(err)
	}
	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	syslog.Info("recv signal and exit:", <-ch)
	service.Stop()
	return
}
