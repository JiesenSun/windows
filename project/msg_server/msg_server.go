package msg_server

import (
	"encoding/json"
	"fmt"
	"time"

	"project/common"
	"project/common/errors"
	"project/common/nsq"
	"project/common/syslog"
	"project/common/util"
	"project/model/message"
	//"project/model/user"
)

type Handler struct {
}

type SendUserMsg struct {
	Type      uint16 `json:"type"`
	ToUid     uint64 `json:"touid"`
	Content   string `json:"content"`
	ExtraData string `json:"extradata"`
}

func (h *Handler) SendMsg(head *common.PackageHead, jsonBody []byte, tail *common.PackageTail) ([]byte, error) {
	var req SendUserMsg
	if err := json.Unmarshal(jsonBody, &req); err != nil {
		return nil, errors.As(err, string(jsonBody))
	}

	var msg message.UserMsgItem
	msg.SendTime = uint32(time.Now().Unix())
	msg.MsgId = message.GetUserMsgID()
	msg.FromUid = tail.UID
	msg.ToUid = req.ToUid
	msg.Content = req.Content
	msg.Type = req.Type
	msg.ExtraData = req.ExtraData

	msgBody, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.As(err, msg)
	}

	msgBuf := common.GetBuffer()
	defer common.PutBuffer(msgBuf)
	head.PkgLen = uint16(common.DATA_PACKAGE_HEAD_SIZE + len(jsonBody))
	if err := common.Package(msgBuf, head, msgBody, tail); err != nil {
		return nil, errors.As(err)
	}
	common.NsqPublish(common.MSG_CENTER_TOPIC, msgBuf.Bytes())
	message.SaveUserMsg(&msg)

	respStr := fmt.Sprintf(`{"msgid":%d,"sendtime":%d}`, msg.MsgId, msg.SendTime)
	return []byte(respStr), nil
}

type RecvUserMsg struct {
	MsgId uint64 `json:"msgid,omitempty"`
}

func (h *Handler) ReceivedUserMsg(head common.PackageHead, jsonBody []byte, tail *common.PackageTail) ([]byte, error) {
	var req RecvUserMsg
	if err := json.Unmarshal(jsonBody, &req); err != nil {
		return nil, errors.As(err, string(jsonBody))
	}

	message.DelUserMsg(tail.UID, req.MsgId)
	return nil, nil
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

	var err error = nil
	resp := []byte("")
	switch head.Command {
	case common.XXX_CMD_SEND_USER_MSG:
		resp, err = h.SendMsg(head, jsonStr, tail)
	case common.XXX_CMD_RECV_USER_MSG:
		//resp, err = h.ReceivedUserMsg(head, jsonStr, &tail)
		common.NsqPublish(common.MSG_STORE_TOPIC, message.Body)
	default:
		syslog.Info("invalid cmd:", head.Command, *tail)
		return nil
	}

	if err != nil {
		syslog.Warn("msg_server msg handle failed:", err, head, tail)
		head.ErrorCode = uint32(errors.Code(err))
		resp = []byte("")
	}
	head.PkgLen = uint16(common.DATA_PACKAGE_HEAD_SIZE + len(resp))
	respBuf := common.GetBuffer()
	defer common.PutBuffer(respBuf)
	if err := common.Package(respBuf, head, resp, tail); err != nil {
		syslog.Info(err)
		return nil
	}
	topic := fmt.Sprintf("conn_%s_%d", util.IPToStr(tail.IP), tail.Port)
	common.NsqPublish(topic, respBuf.Bytes())
	return nil
}

func StartServer() {
	handler := &Handler{}

	for i := 0; i < 10; i++ {
		if _, err := common.NsqConsumerGO(common.MSG_SERVER_TOPIC, "msg_server_channel", uint(10), handler); err != nil {
			panic(err)
		}
	}
	select {}
	syslog.SysLogDeinit()
	return
}
