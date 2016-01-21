package user

import (
	"fmt"
)

func SetUserToken(uid uint64, token string) error {
	return g_redis_client.Set(fmt.Sprint(USER_TOKEN, uid), token, 0).Err()
}

func GetUserToken(uid uint64) (string, error) {
	return g_redis_client.Get(fmt.Sprint(USER_TOKEN, uid)).Result()
}
