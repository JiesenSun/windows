package message

import (
	"time"
)

var (
	MSG_DB                = "dudb"
	MSG_USER_MSG_TABLE    = "user_msg"
	MSG_TEAM_INFO_TABLE   = "team_info"
	MSG_TEAM_MEMBER_TABLE = "team_table"

	KEY_TEAMMSGBUF = "TEAMMSGBUF_"
	SET_TEAMMSGID  = "TEAMMSGID_"
)

const (
	MSG_TYPE_USER = iota
	MSG_TYPE_TEAM
	MSG_TYPE_SYSTEM
)

const (
	MAX_TEAM_MSG_PER = 300
)

func init() {
	if err := UserMsgInit(); err != nil {
		panic(err.Error())
	}
}

/*
* msg id
 */
const (
	MSG_ID_MASK = 0xffffffffffffff
)

func GetUserMsgID() uint64 {
	return uint64(time.Now().UnixNano()) & MSG_ID_MASK
}
