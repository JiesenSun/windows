package user_center

import (
	"project/common/errors"
	"project/common/protocol"
	"project/common/syslog"
	"project/model/user"

	"github.com/golang/protobuf/proto"
)

const (
	LOGIN_POLICY_MANUAL = 0
	LOGIN_POLICY_AUTO   = 1
)

func (h *DBHandler) Login(head *protocol.PackageHead, jsonBody []byte, tail *protocol.PackageTail) (result []byte) {
	var req protocol.LoginRequest
	var rsp protocol.LoginResponse
	errCode := 0
	var err error = nil

	defer func() {
		rsp.ErrorCode = proto.Uint32(uint32(errCode))
		syslog.Debug(rsp.String())
		result, err = proto.Marshal(&rsp)
		if err != nil {
			syslog.Debug(err)
		}
	}()

	if err := proto.Unmarshal(jsonBody, &req); err != nil {
		errCode = errors.Code(errors.CLIENT_ERR_UNKNOW_ERROR)
		return
	}

	phonenum := req.GetUsername()
	password := req.GetPassword()

	userInfo, _ := user.GetUserInfoByPhoneNum(phonenum)
	if userInfo == nil {
		errCode = errors.Code(errors.CLIENT_ERR_USER_NOT_EXIST)
		return
	}

	if userInfo.Password != password {
		errCode = errors.Code(errors.CLIENT_ERR_PASSWORD_ERROR)
		return
	}

	us, err := user.GetUserState(userInfo.Uid)
	if err != nil {
		syslog.Info("GetUserState failed!!!", err)
	}
	us.Uid = userInfo.Uid
	us.Sid = tail.SID
	us.ConnIP = tail.IP
	us.ConnPort = tail.Port
	//us.SetupId = req.SetupId
	us.Online = true
	if err := user.SetUserState(us); err != nil {
		syslog.Info("SetUserState failed", us)
	}

	rsp.UserId = proto.Uint64(userInfo.Uid)
	tail.UID = userInfo.Uid
	return
}
