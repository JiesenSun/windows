package user_center

import (
	"project/common/errors"
	"project/common/protocol"
	"project/common/syslog"
	"project/model/user"

	"github.com/golang/protobuf/proto"
)

func (h *DBHandler) Register(head *protocol.PackageHead, jsonBody []byte, tail *protocol.PackageTail) (result []byte) {
	var req protocol.RegisterRequest
	var rsp protocol.RegisterResponse
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
	if userInfo != nil {
		errCode = errors.Code(errors.CLIENT_ERR_USER_EXIST)
		return
	}

	userInfo = &user.UserInfo{
		Password: password,
		PhoneNum: phonenum,
	}

	if err := user.SaveUserInfo(userInfo); err != nil {
		errCode = errors.Code(errors.CLIENT_ERR_UNKNOW_ERROR)
		return
	}

	rsp.UserId = proto.Uint64(userInfo.Uid)
	tail.UID = userInfo.Uid
	return
}
