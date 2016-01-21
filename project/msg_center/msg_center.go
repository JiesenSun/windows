package msg_center

import (
	"fmt"

	"project/common"
	"project/common/nsq"
	"project/common/syslog"
	"project/common/util"
	"project/model/user"
)

type Handler struct {
}

func (h *Handler) HandleMessage(message *nsq.Message) error {
	dataPackage := common.GetDataPackage()
	defer common.PutDataPackage(dataPackage)
	if err := dataPackage.Unpackage(message.Body); err != nil {
		syslog.Warn(err)
		return nil
	}

	head := &dataPackage.Head
	jsonStr := dataPackage.BodyData()
	tail := &dataPackage.Tail

	switch head.Command {
	case common.XXX_CMD_SEND_USER_MSG:
	default:
		syslog.Warn("Invalid Cmd :", head.Command, *tail)
	}
	// get user state
	us, err := user.GetUserState(tail.UID)
	if err != nil {
		syslog.Info("GetUserState failed!!!", tail)
		return nil
	}
	// not online
	if us.Online == false {
		syslog.Debug("client not online!!! uid:", tail.UID)
		return nil
	}

	tail.IP = us.ConnIP
	tail.Port = us.ConnPort
	tail.SID = us.Sid
	tail.UID = us.Uid

	respBuf := common.GetBuffer()
	defer common.PutBuffer(respBuf)
	if err := common.Package(respBuf, head, jsonStr, tail); err != nil {
		syslog.Info(err)
		return nil
	}
	topic := fmt.Sprintf("conn_%s_%d", util.IPToStr(tail.IP), tail.Port)
	common.NsqPublish(topic, respBuf.Bytes())
	return nil
}

func StartServer() {
	for i := 0; i < 10; i++ {
		if _, err := common.NsqConsumerGO(common.MSG_CENTER_TOPIC, "msgcenter-channel", 10, &Handler{}); err != nil {
			panic(err)
		}
	}
	select {}
	syslog.SysLogDeinit()
	return
}
