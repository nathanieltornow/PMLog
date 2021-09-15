package redis

import (
	"github.com/go-redis/redis/v8"
)

type Server struct {
	redisClient *redis.Client
}

func NewServer() (*Server, error) {
	server := new(Server)
	server.redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return server, nil
}
