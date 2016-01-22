package client

import (
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"project/common/protocol"
)

type RegisterReceiver struct{}

func (*RegisterReceiver) Run(data []byte) {
	resp := &protocol.RegisterResponse{}
	if err := proto.Unmarshal(data, resp); err != nil {
		panic(err.Error())
	}
	println(resp.String())
}

func TestRegister(t *testing.T) {
	socket := NewSocket()
	println("------------statt test-----------")
	if !socket.Connect("10.0.2.15", 8881) {
		panic("connect server failed")
	}

	req := &protocol.RegisterRequest{
		Username: proto.String("18702759796"),
		Password: proto.String("wuxiangan"),
	}

	data, err := proto.Marshal(req)
	if err != nil {
		panic(err.Error())
	}
	cmd := 10001

	socket.Send(cmd, data, &RegisterReceiver{})
	time.Sleep(time.Second * 300)
}
