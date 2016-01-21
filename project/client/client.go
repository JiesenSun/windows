package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	g_client       net.Conn = nil
	cmdPath                 = "./tcp_test_data"
	SIZEOF_PKGHEAD          = 12
)

type PkgHead struct {
	PkgLen    uint16 //整个包的长度
	Version   uint16 //协议版本号
	Command   uint32 //命令字
	ErrorCode uint32 //通讯流水 ID
}

type Client struct {
	Socket net.Conn
}

func (this *Client) Send(head *PkgHead, body []byte) error {
	var sendBuf [2048]byte
	buf := bytes.NewBuffer(sendBuf[:])
	buf.Reset()
	head.PkgLen = uint16(SIZEOF_PKGHEAD + len(body))
	binary.Write(buf, binary.BigEndian, head)
	binary.Write(buf, binary.BigEndian, body)
	_, err := this.Socket.Write(buf.Bytes())
	return err
}

func (this *Client) Recv() error {
	var recvBuf [2048]byte
	head := &PkgHead{}
	if _, err := io.ReadFull(this.Socket, recvBuf[:2]); err != nil {
		return err
	}
	pkgLen := binary.BigEndian.Uint16(recvBuf[:2])
	fmt.Println(pkgLen, recvBuf[0], recvBuf[1])
	if _, err := io.ReadFull(this.Socket, recvBuf[2:pkgLen]); err != nil {
		return err
	}

	respReader := bytes.NewReader(recvBuf[:pkgLen])
	if err := binary.Read(respReader, binary.BigEndian, head); err != nil {
		return err
	}
	jsonStr := make([]byte, int(head.PkgLen)-SIZEOF_PKGHEAD)
	if err := binary.Read(respReader, binary.BigEndian, jsonStr); err != nil {
		return err
	}
	fmt.Println("recv head:", head)
	fmt.Println("recv body:", string(jsonStr))
	return nil
}

func (this *Client) Exec(cmd string) {
	cmdInt, err := strconv.ParseInt(cmd, 10, 32)
	if err != nil {
		panic(err)
	}
	pkgHead := &PkgHead{Command: uint32(cmdInt)}
	file := cmdPath + "/" + cmd
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	if err := this.Send(pkgHead, body); err != nil {
		panic(err)
	}
	return
}

func usage() {
	fmt.Println("命令          描述")
	fmt.Println("10001         注册")
	fmt.Println("10107         找回密码")
	fmt.Println("10002         登录")
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

func (this *Client) Run() {
	var cmd string
	go func() {
		for {
			time.Sleep(time.Minute * 2)
			this.Exec("10000")
		}
	}()
	go func() {
		for {
			if err := this.Recv(); err != nil {
				panic(err)
			}
		}
	}()
	for {
		usage()
		fmt.Scanln(&cmd)
		this.Exec(cmd)
	}
}
func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "10.0.2.15:9100")
	if err != nil {
		panic(err)
	}
	socket, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	//fmt.Println("connect success")
	client := &Client{socket}
	// 先登录
	client.Exec("10002")
	if err := client.Recv(); err != nil {
		fmt.Println(err)
		return
	}
	// 执行命令行
	if len(os.Args) > 1 {
		client.Exec(os.Args[1])
		client.Recv()
		return
	}
	// 循环执行命令
	client.Run()
}
