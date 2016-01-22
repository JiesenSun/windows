package user_center

import (
	"encoding/binary"
	"fmt"

	"project/common"
	"project/common/nsq"
	"project/common/protocol"
	"project/common/syslog"
	"project/common/util"
)

type DBHandler struct {
}

func (h *DBHandler) HandleMessage(message *nsq.Message) error {
	var resp []byte
	var err error
	defer func() {
		if err != nil {
			syslog.Warn(err)
		}
	}()
	dataPackage := protocol.GetDataPackage()
	defer protocol.PutDataPackage(dataPackage)

	if err := dataPackage.Unpackage(message.Body); err != nil {
		return nil
	}
	head := &dataPackage.Head
	jsonStr := dataPackage.BodyData()
	tail := &dataPackage.Tail

	syslog.Debug(*head, string(jsonStr), *tail)

	switch head.Command {
	case protocol.CLIENT_CMD_USER_REGISTER:
		resp = h.Register(head, jsonStr, tail)
	case protocol.CLIENT_CMD_USER_LOGIN:
		resp = h.Login(head, jsonStr, tail)
	default:
		syslog.Debug("invalid cmd:", head.Command, *tail)
		return nil
	}

	respBuf := protocol.GetBuffer()
	defer protocol.PutBuffer(respBuf)
	head.PkgLen = uint16(protocol.DATA_PACKAGE_HEAD_SIZE + len(resp))
	binary.Write(respBuf, binary.BigEndian, head)
	binary.Write(respBuf, binary.BigEndian, resp)
	binary.Write(respBuf, binary.BigEndian, tail)
	topic := fmt.Sprintf("conn_%s_%d", util.IPToStr(tail.IP), tail.Port)
	common.NsqPublish(topic, respBuf.Bytes())
	syslog.Debug2("db_server --> conn_server publish:", topic, head, string(resp), tail)
	return nil
}

func StartServer() {
	handler := &DBHandler{}
	if _, err := common.NsqConsumerGO(common.USER_CENTER_TOPIC, "user_center_channel", 5, handler); err != nil {
		panic(err)
	}
	select {}
}
