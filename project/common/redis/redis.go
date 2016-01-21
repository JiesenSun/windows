package redis

import (
	"gopkg.in/redis.v3"
)

var (
	g_redisClient *Redis = nil
)

type Redis struct {
	*redis.Client
}

func NewRedis(addr string) *Redis {
	client := redis.NewClient(&redis.Options{Addr: addr})
	if _, err := client.Ping().Result(); err != nil {
		panic(err)
	}
	return &Redis{client}
}

func Init(addr string) {
	g_redisClient = NewRedis(addr)
}

func Deinit() {
	g_redisClient.Close()
}

func Get() *Redis {
	return g_redisClient
}
