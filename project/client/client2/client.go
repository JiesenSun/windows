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

	"github.com/bitly/go-simplejson"
)

var (
	g_client       net.Conn = nil
	g_uid          uint64   = 0
	g_sid          uint64   = 0
	cmdPath                 = "./tcp_test_data"
	SIZEOF_PKGHEAD          = 24
)

type PkgHead struct {
	PkgLen uint16 //整个包的长度
	Cmd    uint16 //命令字
	Ver    uint16 //协议版本号
	Seq    uint16 //通讯流水 ID
	Sid    uint32 //Session ID
	Uid    uint64 //UID
	Flag   uint32 //包体压缩 & 包体加密
}

func sendPkg(head *PkgHead, body []byte) error {
	buf := new(bytes.Buffer)
	head.PkgLen = uint16(SIZEOF_PKGHEAD + len(body))
	binary.Write(buf, binary.BigEndian, head)
	binary.Write(buf, binary.BigEndian, body)
	g_client.SetDeadline(time.Now().Add(time.Hour))
	if _, err := g_client.Write(buf.Bytes()); err != nil {
		fmt.Println("g_client.Write error")
		return err
	}
	return nil
}

func recvPkg() error {
	g_client.SetDeadline(time.Now().Add(time.Hour))
	head := &PkgHead{}
	pkgLenSlice := make([]byte, 2)
	if _, err := io.ReadFull(g_client, pkgLenSlice); err != nil {
		fmt.Println("io.ReadFull error", err)
		return err
	}
	pkgLenInt := binary.BigEndian.Uint16(pkgLenSlice)
	respBuf := make([]byte, pkgLenInt)
	respBuf[0] = pkgLenSlice[0]
	respBuf[1] = pkgLenSlice[1]

	//g_client.SetDeadline(time.Now().Add(time.Second))
	if _, err := io.ReadFull(g_client, respBuf[2:]); err != nil {
		fmt.Println("io.ReadFull error", err)
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
		println("server return failed:", head.Cmd)
		return nil
	}
	fmt.Println("recv head:", head)
	fmt.Println("recv body:", string(jsonStr))
	//fmt.Println("server return success")
	if head.Cmd == 10102 && len(jsonStr) != 0 {
		js, _ := simplejson.NewJson(jsonStr)
		uid, _ := js.Get("uid").Uint64()
		sid, _ := js.Get("sid").Uint64()
		g_sid = sid
		g_uid = uid
	}
	return nil
}

func exec(cmd string) {
	cmdInt, err := strconv.ParseInt(cmd, 10, 32)
	if err != nil {
		panic(err)
	}
	pkgHead := &PkgHead{
		Cmd: uint16(cmdInt),
		Sid: uint32(g_sid),
		Uid: g_uid,
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
	if err := sendPkg(pkgHead, body); err != nil {
		panic(err)
	}
	return
}

func usage() {
	fmt.Println("命令          描述")
	fmt.Println("10101         注册")
	fmt.Println("10102         登录")
	fmt.Println("30101         发送用户消息")
	fmt.Println("30201         确认用户消息")
	fmt.Println("30102         发送群消息")
	fmt.Println("30202         确认群消息")
	fmt.Println("50001         创建群")
	fmt.Println("50002         删除群")
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

func run() {
	var cmd string
	go func() {
		for {
			time.Sleep(time.Minute * 2)
			exec("10100")
		}
	}()
	go func() {
		for {
			if err := recvPkg(); err != nil {
				panic(err)
			}
		}
	}()
	for {
		usage()
		fmt.Scanln(&cmd)
		exec(cmd)
	}
}
func main() {
	//tcpAddr, err := net.ResolveTCPAddr("tcp", "192.168.20.51:9100")
	//tcpAddr, err := net.ResolveTCPAddr("tcp", "112.74.66.141:9100")
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9100")
	if err != nil {
		panic(err)
	}
	client, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	//fmt.Println("connect success")
	g_client = client
	// 先登录
	exec("10102")
	recvPkg()
	// 执行命令行
	if len(os.Args) > 1 {
		exec(os.Args[1])
		recvPkg()
		return
	}
	// 循环执行命令
	run()

	g_client.Close()
}
