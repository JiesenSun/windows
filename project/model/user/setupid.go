package user

import (
	"fmt"
)

func SetUserSetupId(uid uint64, setupId string) error {
	return g_redis_client.Set(fmt.Sprint(USER_TOKEN, uid), setupId, 0).Err()
}

func GetUserSetupId(uid uint64) (string, error) {
	return g_redis_client.Get(fmt.Sprint(USER_TOKEN, uid)).Result()
}

func IsSameSetupId(uid uint64, setupid string) bool {
	id, _ := GetUserSetupId(uid)
	if id == setupid {
		return true
	}

	return false
}
