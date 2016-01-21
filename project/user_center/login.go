package user_center

import (
	"encoding/json"

	"project/common/errors"
	"project/common/protocol"
	"project/common/syslog"
	"project/model/user"
)

const (
	LOGIN_POLICY_MANUAL = 0
	LOGIN_POLICY_AUTO   = 1
)

type LoginReq struct {
	Platform string `json:"platform,omitempty"`
	Uid      uint64 `json:"uid,omitempty"`
	Password string `json:"password,omitempty"`
	Policy   int    `json:"policy,omitempty"`
	SetupId  string `json:"setupid,omitempty"`
}

type LoginResp struct {
	Uid int64 `json:"uid,omitempty"`
	Sid int   `json:"sid,omitempty"`
}

func (h *DBHandler) Login(head *protocol.PackageHead, jsonBody []byte, tail *protocol.PackageTail) ([]byte, error) {
	var req LoginReq
	if err := json.Unmarshal(jsonBody, &req); err != nil {
		return []byte(""), errors.As(err, string(jsonBody))
	}

	userInfo, err := user.GetUserInfoByUid(req.Uid)
	if err != nil || userInfo == nil {
		return []byte(""), errors.CLIENT_ERR_USER_NOT_EXIST
	}

	if userInfo.Password != req.Password {
		return []byte(""), errors.CLIENT_ERR_PASSWORD_ERROR
	}

	us, err := user.GetUserState(req.Uid)
	if err != nil {
		syslog.Info("GetUserState failed!!!", err)
	}
	us.Uid = req.Uid
	us.Sid = tail.SID
	us.ConnIP = tail.IP
	us.ConnPort = tail.Port
	us.SetupId = req.SetupId
	us.Online = true
	if err := user.SetUserState(us); err != nil {
		syslog.Info("SetUserState failed", us)
	}

	resp := LoginResp{
		Uid: int64(userInfo.Uid),
		Sid: int(tail.SID),
	}
	respBuf, err := json.Marshal(resp)
	if err != nil {
		return []byte(""), errors.As(err)
	}

	tail.UID = req.Uid
	return respBuf, nil
}
