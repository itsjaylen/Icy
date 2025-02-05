package redis

import (
	"context"
	logger "itsjaylen/IcyLogger"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient initializes and returns a Redis client.
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	logger.Info.Println("Connected to Redis successfully")
	return &RedisClient{Client: client}, nil
}
