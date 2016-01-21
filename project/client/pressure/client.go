package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/bitly/go-simplejson"
	_ "sirendaou.com/duserver/common/syscall"
)

var (
	cmdPath        = "./tcp_test_data"
	SIZEOF_PKGHEAD = 24
	MAX_CONN_NUM   = 50000
	MIN_UID        = 1000000
	cur_count      = 0

	g_rwMutex    = &sync.RWMutex{}
	g_client_map = make(map[uint64]*Client)
)

type Client struct {
	conn net.Conn
	uid  uint64
	sid  uint64
}

type PkgHead struct {
	PkgLen uint16 //整个包的长度
	Cmd    uint16 //命令字
	Ver    uint16 //协议版本号
	Seq    uint16 //通讯流水 ID
	Sid    uint32 //Session ID
	Uid    uint64 //UID
	Flag   uint32 //包体压缩 & 包体加密
}

func (c *Client) sendPkg(head *PkgHead, body []byte) error {
	var seq uint16 = 1
	buf := new(bytes.Buffer)
	head.Seq = seq
	seq++
	head.PkgLen = uint16(SIZEOF_PKGHEAD + len(body))
	binary.Write(buf, binary.BigEndian, head)
	binary.Write(buf, binary.BigEndian, body)
	c.conn.SetDeadline(time.Now().Add(time.Hour))
	if _, err := c.conn.Write(buf.Bytes()); err != nil {
		fmt.Println("c.conn.Write error")
		return err
	}
	return nil
}

func (c *Client) recvPkg() error {
	c.conn.SetDeadline(time.Now().Add(time.Hour))
	head := &PkgHead{}
	pkgLenSlice := make([]byte, 2)
	if _, err := io.ReadFull(c.conn, pkgLenSlice); err != nil {
		log.Println("io.ReadFull error", err, *c, c.conn.LocalAddr().String())
		return err
	}
	pkgLenInt := binary.BigEndian.Uint16(pkgLenSlice)
	respBuf := make([]byte, pkgLenInt)
	respBuf[0] = pkgLenSlice[0]
	respBuf[1] = pkgLenSlice[1]

	//c.conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.ReadFull(c.conn, respBuf[2:]); err != nil {
		log.Println("io.ReadFull error", err, *c, c.conn.LocalAddr().String())
		return err
	}

	respReader := bytes.NewReader(respBuf)
	if err := binary.Read(respReader, binary.BigEndian, head); err != nil {
		fmt.Println("binary.Read error:")
		return err
	}
	jsonStr := make([]byte, int(head.PkgLen)-SIZEOF_PKGHEAD)
	if err := binary.Read(respReader, binary.BigEndian, jsonStr); err != nil {
		fmt.Println("binary.Read error:")
		return err
	}
	if head.Sid != 0 {
		println("server return failed:", head.Cmd, head.Sid)
		return nil
	}
	fmt.Println("recv head:", head)
	fmt.Println("recv body:", string(jsonStr))
	c.recvDeal(head, jsonStr)
	return nil
}
func (c *Client) recvDeal(head *PkgHead, jsonStr []byte) {
	js, _ := simplejson.NewJson(jsonStr)
	//fmt.Println("server return success")
	if head.Cmd == 10102 {
		uid, _ := js.Get("uid").Uint64()
		sid, _ := js.Get("sid").Uint64()
		c.sid = sid
		c.uid = uid
		fmt.Println("login success")
		g_rwMutex.Lock()
		g_client_map[c.uid] = c
		g_rwMutex.Unlock()
	} else if head.Cmd == 35101 {
		msgid, _ := js.Get("msgid").Uint64()
		c.execArgs("30201", c.uid, msgid)
	}
}

func (c *Client) execArgs(cmd string, args ...interface{}) {
	cmdInt, err := strconv.ParseInt(cmd, 10, 32)
	if err != nil {
		panic(err)
	}
	pkgHead := &PkgHead{
		Cmd: uint16(cmdInt),
		Sid: uint32(c.sid),
		Uid: c.uid,
	}
	bodyStr := ""
	switch cmd {
	case "30201":
		bodyStr = fmt.Sprintf(`{"uid":%v,"msgid":%v}`, args...)
	case "10102":
		bodyStr = fmt.Sprintf(`{"platform": "i", "uid": %v, "password": "%v"}`, args...)
	case "30101":
		bodyStr = fmt.Sprintf(`{"msgcontent":"testSend", "apnstext":"test", "touid":%v}`, args...)
	}
	fmt.Println(cmd, bodyStr)
	if err := c.sendPkg(pkgHead, []byte(bodyStr)); err != nil {
		panic(err)
	}
}

func (c *Client) exec(cmd string) {
	cmdInt, err := strconv.ParseInt(cmd, 10, 32)
	if err != nil {
		panic(err)
	}
	pkgHead := &PkgHead{
		Cmd: uint16(cmdInt),
		Sid: uint32(c.sid),
		Uid: c.uid,
	}
	file := cmdPath + "/" + cmd
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("open file failed:", err)
		return
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	if err := c.sendPkg(pkgHead, body); err != nil {
		panic(err)
	}
	return
}

func (c *Client) run() {
	//var cmd string
	exit := false
	go func() {
		for !exit {
			time.Sleep(time.Minute * 2)
			c.execArgs("10100")
		}
	}()
	go func() {
		for !exit {
			err := c.recvPkg()
			if err == nil {
				continue
			} else if errno, ok := err.(syscall.Errno); ok && errno == syscall.ECONNRESET {
				exit = true
				log.Println(err, *c)
				return
				go CreateClient(c.uid, "test")
			} else {
				panic(err)
			}
		}
	}()
	/*
		for {
			//usage()
			fmt.Scanln(&cmd)
			c.exec(cmd)
		}
	*/
}

func (c *Client) Close() {
	c.conn.Close()
}

func CreateClient(uid uint64, password string) {
RECONNECTION:
	//tcpAddr, err := net.ResolveTCPAddr("tcp", "192.168.20.51:9100")
	tcpAddr, err := net.ResolveTCPAddr("tcp", "112.74.66.141:9100")
	//tcpAddr, err := net.ResolveTCPAddr("tcp", "112.74.66.141:9100")
	//tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9100")
	if err != nil {
		fmt.Println(err)
		return
	}
	clientIP := fmt.Sprintf("127.%v.%v.%v:0", rand.Uint32()%255, rand.Uint32()%255, rand.Uint32()%255)
	locAddr, err := net.ResolveTCPAddr("tcp", clientIP)
	if err != nil {
		panic(err)
	}
	connClient, err := net.DialTCP("tcp", locAddr, tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &Client{
		conn: connClient,
		uid:  uid,
	}
	log.Println("connect start:", *client)
	err = client.recvPkg()
	if err != nil && err != io.EOF {
		log.Println(err, *client)
		return
		goto RECONNECTION
	} else if err != nil {
		panic(err)
	}
	log.Println("connect end:", *client)
	// 先登录
	client.execArgs("10102", uid, password)
	client.recvPkg()
	cur_count++
	go client.run()

	//select {}
	time.Sleep(3 * time.Minute)
	for {
		// 跑业务
		//second := time.Duration(rand.Uint32()%150) + 30
		time.Sleep(time.Second * 60)
		if client.sid == 0 {
			continue
		}
		uid := int(rand.Uint32())%cur_count + MIN_UID
		g_rwMutex.RLock()
		_, ok := g_client_map[uint64(uid)]
		g_rwMutex.RUnlock()
		if ok {
			client.execArgs("30101", uid)
		}
	}
	client.Close()
}

func usage() {
	fmt.Println("命令          描述")
	fmt.Println("10101         注册")
	fmt.Println("10107         找回密码")
	fmt.Println("10102         登录")
	fmt.Println("30101         发送用户消息")
	fmt.Println("30201         确认用户消息")
	fmt.Println("30102         发送群消息")
	fmt.Println("30202         确认群消息")
	fmt.Println("50001         创建群")
	fmt.Println("50002         删除群")
	fmt.Println("50003         搜索群")
	fmt.Println("50011         获得群信息")
	fmt.Println("50012         获得用户群列表")
	fmt.Println("50021         修改群消息")
	fmt.Println("50022         添加群成员")
	fmt.Println("50023         删除群成员")
	fmt.Println("50024         获取群成员")
	fmt.Println("50025         添加好友(1 - 白 2 - 黑)")
	fmt.Println("50026         删除好友(1 - 白 2 - 黑)")
	fmt.Println("50027         查询好友(1 - 白 2 - 黑)")
}

func main() {
	if len(os.Args) > 1 {
		min_uid, _ := strconv.ParseInt(os.Args[1], 10, 32)
		MIN_UID = int(min_uid)
	}
	for i := 0; i < MAX_CONN_NUM; i++ {
		time.Sleep(time.Millisecond * 10)
		go CreateClient(uint64(i+MIN_UID), "test")
	}
	fmt.Println("finish")
	select {}
}
