package message

import (
	"fmt"
	"testing"
	"time"

	"sirendaou.com/duserver/common"
)

func TestUserMsg(t *testing.T) {
	common.MongoInit("127.0.0.1")
	if err := UserMsgInit(); err != nil {
		fmt.Println(err)
		return
	}
	msg := &UserMsgItem{
		MsgId:   1,
		FromUid: 2,
		ToUid:   3,
		Content: "hello world",
	}
	if err := SaveUserMsg(msg); err != nil {
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second)
	msgs, _ := GetUserMsgList(msg.ToUid)
	DelUserMsg(msg.ToUid, msg.MsgId)
	for _, m := range msgs {
		fmt.Println(*m)
	}
	fmt.Println("msg count:", len(msgs))
	time.Sleep(4 * time.Second)
}
