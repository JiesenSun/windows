package user

import (
	"sync"

	"project/common/redis"
)

var (
	USER_DB           = "dudb"
	USER_FRIEND_TABLE = "friend_list"
	USER_INFO_TABLE   = "user_info"

	SET_WHITELIST = "WHITE_"
	SET_BLACKLIST = "BLACK_"
	USER_TOKEN    = "redis_token_"

	MIN_USER_ID     = 100000
	user_mutex_lock = &sync.Mutex{}

	g_redis_client = redis.Get()
)

func init() {
	if err := FriendInit(); err != nil {
		panic(err.Error())
	}
	if err := UserInfoInit(); err != nil {
		panic(err.Error())
	}
}
