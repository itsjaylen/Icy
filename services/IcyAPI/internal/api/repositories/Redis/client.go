package redis

import (
	"IcyAPI/internal/utils"
	"context"
	"time"

	logger "itsjaylen/IcyLogger"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis client instance.
type RedisClient struct {
	Client   *redis.Client
	Addr     string
	Password string
	DB       int
}

// NewRedisClient initializes and returns a Redis client with retry logic.
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := &RedisClient{
		Addr:     addr,
		Password: password,
		DB:       db,
	}

	err := utils.Retry(5, 2*time.Second, client.connect) // Use the retry function
	if err != nil {
		return nil, err
	}
	return client, nil
}

// connect establishes a connection to Redis and performs a health check.
func (r *RedisClient) connect() error {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})

	if err := r.Client.Ping(context.Background()).Err(); err != nil {
		return err
	}

	logger.Info.Println("Connected to Redis successfully")
	return nil
}

// Reconnect attempts to reconnect to Redis using the retry utility.
func (r *RedisClient) Reconnect() {
	err := utils.Retry(5, 2*time.Second, r.connect)
	if err != nil {
		logger.Error.Println("Failed to reconnect to Redis after multiple attempts")
	} else {
		logger.Info.Println("Reconnected to Redis successfully")
	}
}

// Close gracefully closes the Redis client connection.
func (r *RedisClient) Close() {
	if err := r.Client.Close(); err != nil {
		logger.Warn.Println("Error closing Redis connection:", err)
	} else {
		logger.Info.Println("Redis connection closed")
	}
}
