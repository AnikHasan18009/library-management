package Redis

import (
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetRedisClient() *redis.Client {
	return redisClient
}
