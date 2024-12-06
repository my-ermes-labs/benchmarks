package mobile_client

import (
	rc "github.com/ermes-labs/storage-redis/packages/go"
	"github.com/redis/go-redis/v9"
)

var db = 10

type Commands struct {
	rc.RedisCommands
}

func NewCommands() *Commands {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   db,
	})

	db--

	return &Commands{
		RedisCommands: *rc.NewRedisCommands(redisClient),
	}
}
