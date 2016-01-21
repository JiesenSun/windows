package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"project/common/protocol"
	"project/common/safemap"
	"project/common/sync"
)

const (
	g_socket_timeout = 5 * time.Second
)

type Receiver interface {
	Run([]byte)
}

type Socket struct {
	svrAddr        *net.TCPAddr
	socket         *net.TCPConn
	sendQueue      chan []byte
	waitGroup      *sync.WaitGroup
	recvHandle     *safemap.Map
	callbackHandle *safemap.Map
	connectChan    chan struct{}
	pkgSeq         uint32
}

func NewSocket() *Socket {
	socket := &Socket{
		svrAddr:        nil,
		socket:         nil,
		sendQueue:      make(chan []byte, 1000),
		connectChan:    make(chan struct{}, 1),
		recvHandle:     safemap.NewMap(),
		callbackHandle: safemap.NewMap(),
		waitGroup:      sync.NewWaitGroup(),
	}
	go socket.send()
	go socket.recv()
	go socket.heartbeat()
	return socket
}

func (s *Socket) connectLock() bool {
	select {
	case s.connectChan <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Socket) connectUnlock() {
	<-s.connectChan
}

func (s *Socket) recv() {
	s.waitGroup.AddOne()
	defer s.waitGroup.Done()

	var dataHead protocol.PackageHead
	var dataLen uint16
	var buf [protocol.DATA_PACKAGE_MAX_SIZE]byte
	byteBuf := bytes.NewBuffer(nil)
	exitNotify := s.waitGroup.ExitNotify()
	for {
		if s.socket == nil {
			time.Sleep(time.Second)
			s.reconnect()
			continue
		}
		select {
		case <-exitNotify:
			return
		default:
		}
		byteBuf.Reset()
		s.socket.SetDeadline(time.Now().Add(g_socket_timeout))
		if _, err := io.ReadFull(s.socket, buf[:2]); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			s.reconnect()
			g_logger.Debug(err.Error())
			continue
		}
		byteBuf.Write(buf[:2])
		if err := binary.Read(byteBuf, binary.BigEndian, &dataLen); err != nil {
			g_logger.Info(err.Error())
		}
		if dataLen < protocol.DATA_PACKAGE_MIN_SIZE || dataLen > protocol.DATA_PACKAGE_MAX_SIZE {
			s.reconnect()
			g_logger.Debug("data package error, data len:" + fmt.Sprint(dataLen))
			continue
		}
		s.socket.SetDeadline(time.Now().Add(g_socket_timeout))
		if _, err := io.ReadFull(s.socket, buf[2:dataLen]); err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			// reconnect
			s.reconnect()
			g_logger.Debug(err.Error())
			continue
		}
		g_logger.Debug("success receive response message!!!")

		byteBuf.Write(buf[:dataLen])
		binary.Read(byteBuf, binary.BigEndian, &dataHead)
		dataLen -= protocol.DATA_PACKAGE_HEAD_SIZE
		binary.Read(byteBuf, binary.BigEndian, buf[:dataLen])
		// handle recv data
		receiver, ok := s.callbackHandle.Get(dataHead.Command).(Receiver)
		if ok && receiver != nil {
			receiver.Run(buf[:dataLen])
			continue
		}
		receiver, ok = s.recvHandle.Delete(dataHead.SequenceID).(Receiver)
		if ok && receiver != nil {
			receiver.Run(buf[:dataLen])
			continue
		}
	}
}

func (s *Socket) send() {
	s.waitGroup.AddOne()
	defer s.waitGroup.Done()

	var totalLen int
	var sendLen int
	var data []byte
	exitNotify := s.waitGroup.ExitNotify()
	for {
		if s.socket == nil {
			time.Sleep(time.Second)
			s.reconnect()
			continue
		}
		select {
		case <-exitNotify:
			return
		case data = <-s.sendQueue:
		}
		totalLen = len(data)
		sendLen = 0
		for sendLen < totalLen {
			s.socket.SetDeadline(time.Now().Add(g_socket_timeout))
			n, err := s.socket.Write(data)
			if err != nil {
				s.reconnect()
				s.sendQueue <- data
				break
			}
			sendLen += n
		}
	}
}

func (s *Socket) Connect(ip string, port int) bool {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return false
	}
	if false == s.connectLock() {
		return false
	}
	defer s.connectUnlock()
	socket, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return false
	}
	s.svrAddr = addr
	s.socket = socket
	return true
}

func (s *Socket) reconnect() {
	if s.svrAddr == nil {
		return
	}

	if false == s.connectLock() {
		return
	}
	defer s.connectUnlock()

	if s.socket != nil {
		s.socket.Close()
		s.socket = nil
	}

	socket, err := net.DialTCP("tcp", nil, s.svrAddr)
	if err != nil {
		g_logger.Debug("reconnect server failed!!! " + err.Error())
	}
	s.socket = socket
}

func (s *Socket) IsConnect() bool {
	return s.socket != nil
}

func (s *Socket) Close() {
	s.waitGroup.Wait()
	s.socket.Close()
	s.socket = nil
}

func (s *Socket) Send(cmd int, data []byte, f Receiver) {
	var head protocol.PackageHead
	head.Command = uint32(cmd)
	head.PkgLen = uint16(protocol.DATA_PACKAGE_HEAD_SIZE + len(data))
	head.SequenceID = s.pkgSeq

	buf := bytes.NewBuffer(make([]byte, head.PkgLen))
	binary.Write(buf, binary.BigEndian, &head)
	binary.Write(buf, binary.BigEndian, data)
	s.sendQueue <- buf.Bytes()
	// 设置接收处理
	s.recvHandle.Set(s.pkgSeq, f)
	// 包序自增
	s.pkgSeq++
}

func (s *Socket) RegisterCallback(cmd int, f Receiver) {
	s.callbackHandle.Set(cmd, f)
}

func (s *Socket) heartbeat() {
	s.waitGroup.AddOne()
	defer s.waitGroup.Done()

	var head protocol.PackageHead
	head.Command = protocol.CLIENT_CMD_HEARTBEAT
	head.PkgLen = uint16(protocol.DATA_PACKAGE_HEAD_SIZE)
	tick := time.Tick(time.Minute * 2)
	exit := s.waitGroup.ExitNotify()

	for {
		select {
		case <-tick:
		case <-exit:
			return
		}
		if s.socket != nil {
			binary.Write(s.socket, binary.BigEndian, &head)
		}
	}
}
