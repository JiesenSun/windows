package user

import (
	"fmt"
	"testing"
)

func TestUserState(t *testing.T) {
	fmt.Println("================")
	state := &UserState{Uid: 3}

	if err := g_redis_client.Set("user_state", state, 0).Err(); err != nil {
		panic(err)
		return
	}

	val, err := g_redis_client.Get("user_state").Result()
	if err != nil {
		panic(err)
		return
	}
	fmt.Println(val)
}
