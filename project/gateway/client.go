package gateway

import (
	"encoding/binary"
	"io"
	"net"
	"sync/atomic"
	"time"

	"project/common"
	"project/common/errors"
	"project/common/pool"
	"project/common/protocol"
	"project/common/syslog"
	"project/model/user"
)

type Client struct {
	IP          uint32
	Port        uint32
	Conn        net.Conn
	Index       int
	Uid         uint64
	Sid         uint32
	Stat        uint32
	LastTime    int64
	MsgCount    int32
	RespCh      chan *protocol.DataPackage
	Service     *Service
	ClientGroup *ClientGroup
}

var (
	CLIENT_TYPE = "CLIENT"
)

func init() {
	pool.RegisterType(CLIENT_TYPE, Client{}, PROGRAM_SCALE)
}

func GetClient() *Client {
	return pool.Get(CLIENT_TYPE).(*Client)
}

func PutClient(client *Client) {
	pool.Put(CLIENT_TYPE, client)
}

func (c *Client) Hello() {
	//dataPackage := pool.Get(DATA_PACKAGE).(*protocol.DataPackage)
	dataPackage := protocol.GetDataPackage()
	dataPackage.Head.PkgLen = uint16(protocol.DATA_PACKAGE_HEAD_SIZE)
	dataPackage.Head.Command = protocol.CLIENT_CMD_HEARTBEAT
	buf := protocol.GetBuffer()
	binary.Write(buf, binary.BigEndian, dataPackage.Head)
	c.Conn.SetDeadline(time.Now().Add(NET_WRITE_DATELINE))
	if _, err := c.Conn.Write(buf.Bytes()); err != nil {
		syslog.Warn(err)
	}
	protocol.PutBuffer(buf)
	protocol.PutDataPackage(dataPackage)
}

type ClientGroup struct {
	Exit     bool
	Service  *Service
	ClientCh chan *Client
	Client   [CLIENT_NUM_PER_GROUP]*Client
	Next     [CLIENT_NUM_PER_GROUP]int
	FreeHead int
	UsedHead int
	Size     int
	RespCh   chan *protocol.DataPackage
}

func NewClientGroup(service *Service) *ClientGroup {
	clientGrp := &ClientGroup{
		Exit:     false,
		Service:  service,
		ClientCh: service.ClientCh,
		UsedHead: CLIENT_NUM_PER_GROUP,
		RespCh:   make(chan *protocol.DataPackage, 1000),
	}
	for i := 0; i < CLIENT_NUM_PER_GROUP; i++ {
		clientGrp.Next[i] = i + 1
	}
	go clientGrp.Request()
	go clientGrp.Response()
	atomic.AddInt32(&service.ClientGroupCount, 1)

	return clientGrp
}

func (this *ClientGroup) Stop() {
	atomic.AddInt32(&this.Service.ClientGroupCount, -1)
	this.Exit = true
}

func (this *ClientGroup) Request() {
	var index int = 0
	var next, prev int
	var client *Client = nil
	for this.Exit == false {
		client = nil
		// 是否存在空闲索引，即组不满
		if this.Size != CLIENT_NUM_PER_GROUP {
			select {
			case client = <-this.ClientCh:
			default:
			}
		}

		if client == nil && this.Size == 0 {
			clientCount := atomic.LoadInt32(&this.Service.ClientCount)
			groupCount := atomic.LoadInt32(&this.Service.ClientGroupCount)
			// 预留10个组，当客户端数量过少时减少组
			if clientCount < ((groupCount-10)*CLIENT_NUM_PER_GROUP) && groupCount > CLIENT_GROUP_NUM {
				this.Stop()
				return
			}
			// 阻塞等待客户端
			select {
			case client = <-this.ClientCh:
			}
		}

		if client != nil { // 有新连接到来，添加组中
			client.Index = this.FreeHead
			this.Client[this.FreeHead] = client
			// 空闲链表
			this.FreeHead = this.Next[this.FreeHead]
			// 非空闲链表
			this.Next[client.Index] = this.UsedHead
			this.UsedHead = client.Index
			this.Size++

			client.RespCh = this.RespCh
			client.ClientGroup = this
			// 确认链接建立
			/*
				* 经测试发现当客户端链接比较多和频繁时，服务器可能处3次握手的第3阶段，
				×服务端连接未完成，客户端完成第2阶段并完成链接，此时发送数据导致服务
				×端在第3阶段关闭连接，造成链接失败，故让服务先发送数据，确保链接完成
				* client.Hello()
			*/
			//client.Hello()
		}
		// 读已连接的客户端请求
		prev = CLIENT_NUM_PER_GROUP
		index = this.UsedHead
		for index != CLIENT_NUM_PER_GROUP {
			next = this.Next[index]
			if ok := this.Client[index].Request(); ok == false {
				this.Client[index].Index = -1
				// 空闲链表
				this.Next[index] = this.FreeHead
				this.FreeHead = index
				// 非空闲链表
				if index == this.UsedHead {
					this.UsedHead = next
				} else {
					this.Next[prev] = next
				}
				this.Client[index].Stop()
				this.Client[index] = nil
				this.Size--

				index = next
				continue
			}
			prev = index
			index = next
		}
	}
}

func (this *ClientGroup) Response() {
	var dp *protocol.DataPackage = nil
	var client *Client = nil
	for this.Exit == false {
		select {
		case dp = <-this.RespCh:
			client = dp.Data.(*Client)
			if err := client.Response(dp); err != nil {
				syslog.Warn(err)
			}
			atomic.AddInt32(&client.MsgCount, -1)
			protocol.PutDataPackage(dp)
		}
	}
}

func (c *Client) Stop() {
	if 0 < atomic.LoadInt32(&c.MsgCount) {
		time.AfterFunc(time.Second*30, c.Stop)
		return
	}
	syslog.Debug("client ", c.Uid, "Sid ", c.Sid, " disconnect!!", c.Conn.RemoteAddr())
	c.Conn.Close()
	us, err := user.GetUserState(c.Uid)
	if err != nil {
		syslog.Info("GetUserState failed!!!", err)
	}
	if us.Uid == c.Uid && us.Sid == c.Sid {
		if err := user.SetUserState(&user.UserState{Uid: c.Uid}); err != nil {
			syslog.Info("SetUserState failed!!!", err)
		}
	}
	c.Stat = CLIENT_STATE_OFFLINE
	c.Service.RWMutex.Lock()
	delete(c.Service.ClientMap, c.Sid)
	c.Service.RWMutex.Unlock()
	atomic.AddInt32(&c.Service.ClientCount, -1)
	PutClient(c)
}

func (c *Client) Request() bool {
	var reqBuf [2048]byte
	curTime := time.Now().Unix()
	c.Conn.SetDeadline(time.Now().Add(NET_READ_DATELINE))
	if _, err := io.ReadFull(c.Conn, reqBuf[:2]); err != nil {
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			// 太长时间没有交互退出
			if c.Stat == CLIENT_STATE_INIT && c.LastTime+CONNECT_NO_LOGIN_TIMEOUT < curTime {
				syslog.Info(c.Uid, c.Sid, "no login timeout 30 s, lastTime:", c.LastTime, " nowTime:", curTime)
				return false
			} else if c.LastTime+CONNECT_NO_REQUEST_TIMEOUT < curTime {
				syslog.Info(c.Uid, c.Sid, "timeout 180 s, lastTime:", c.LastTime, " nowTime:", curTime)
				return false
			}
			return true
		}
		syslog.Info(c.Uid, c.Sid, "recv error:", err.Error())
		return false
	} else {
		headLen := binary.BigEndian.Uint16(reqBuf[0:2])
		if headLen > protocol.DATA_PACKAGE_MAX_SIZE {
			syslog.Info(c.Uid, c.Sid, "head len :", headLen, " too large", reqBuf[0], reqBuf[1])
			return true
		} else if headLen < protocol.DATA_PACKAGE_HEAD_SIZE {
			syslog.Info(c.Uid, c.Sid, "head len :", headLen, " too small", reqBuf[0], reqBuf[1])
			return true
		}
		c.Conn.SetDeadline(time.Now().Add(NET_READ_MAX_DATELINE))
		if _, err := io.ReadFull(c.Conn, reqBuf[2:headLen]); err != nil {
			syslog.Info(c.Uid, c.Sid, "recv error:", err)
			return false
		}
		dataPackage := protocol.GetDataPackage()
		defer protocol.PutDataPackage(dataPackage)
		if err := dataPackage.Unpackage(reqBuf[:headLen]); err != nil {
			syslog.Info(c.Uid, c.Sid, "DecPkgBody failed :", err)
			return false
		}
		head := dataPackage.Head
		syslog.Debug("client requst:", head)
		if head.Command != protocol.CLIENT_CMD_HEARTBEAT && head.Command != protocol.CLIENT_CMD_USER_LOGIN &&
			head.Command != protocol.CLIENT_CMD_USER_REGISTER && c.Stat != CLIENT_STATE_ONLINE {
			syslog.Info("no login no operate!!!", head)
			return false
		} else {
			c.LastTime = time.Now().Unix()
			dataPackage.Tail.IP = g_IP
			dataPackage.Tail.Port = g_Port
			dataPackage.Tail.SID = c.Sid
			dataPackage.Tail.UID = c.Uid
			dataPackage.Data = c
			if head.Command != protocol.CLIENT_CMD_HEARTBEAT {
				//派发消息到其它服务器处理
				err = c.Produce(dataPackage)
			} else {
				//回复心跳
				err = c.Response(dataPackage)
			}
			if err != nil {
				syslog.Info(err)
				return false
			}
		}
	}
	return true
}

func (c *Client) Response(msg *protocol.DataPackage) error {
	head := msg.Head
	tail := msg.Tail

	if (head.Command == protocol.CLIENT_CMD_USER_LOGIN || head.Command == protocol.CLIENT_CMD_USER_REGISTER) && tail.UID != 0 {
		c.Uid = tail.UID
		c.Stat = CLIENT_STATE_ONLINE
	}

	buf := protocol.GetBuffer()
	defer protocol.PutBuffer(buf)
	binary.Write(buf, binary.BigEndian, head)
	binary.Write(buf, binary.BigEndian, msg.BodyData())

	// write to Conn
	c.Conn.SetDeadline(time.Now().Add(NET_WRITE_DATELINE))
	if _, err := c.Conn.Write(buf.Bytes()); err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			syslog.Info("write timeouat:", err, *c)
			return nil
		}
		syslog.Warn("write to conn fail :", err.Error(), msg.Tail)
		return errors.As(err)
	} else {
		syslog.Debug2("write to conn success:", head, string(msg.BodyData()), msg.Tail)
	}
	return nil
}

func (c *Client) Produce(msg *protocol.DataPackage) error {
	head := msg.Head
	buf := protocol.GetBuffer()
	defer protocol.PutBuffer(buf)
	if err := msg.Package(buf); err != nil {
		return errors.As(err)
	}

	if head.Command == protocol.CLIENT_CMD_SEND_USER_MSG || head.Command == protocol.CLIENT_CMD_RECV_USER_MSG {
		if err := common.NsqPublish(common.MSG_SERVER_TOPIC, buf.Bytes()); err != nil {
			return errors.As(err)
		}
		syslog.Debug("gateway --> msg_server publish message:", head, msg.Tail)
	} else if head.Command == protocol.CLIENT_CMD_USER_LOGIN || head.Command == protocol.CLIENT_CMD_USER_REGISTER {
		if err := common.NsqPublish(common.USER_CENTER_TOPIC, buf.Bytes()); err != nil {
			return errors.As(err)
		}
		syslog.Debug("gateway --> user_center publish message:", head, msg.Tail)
	}
	return nil
}
