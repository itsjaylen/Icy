// Package redis provides a Redis client with retry logic.
package redis

import (
	"context"
	"time"

	"github.com/itsjaylen/IcyAPI/internal/utils"
	"github.com/redis/go-redis/v9"
	logger "itsjaylen/IcyLogger"
)

// Client wraps the Redis client instance.
type Client struct {
	Client   *redis.Client
	Addr     string
	Password string
	DB       int
}

// NewClient initializes and returns a Redis client with retry logic.
func NewClient(addr, password string, db int) (*Client, error) {
	client := &Client{
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
func (rd *Client) connect() error {
	rd.Client = redis.NewClient(&redis.Options{
		Addr:     rd.Addr,
		Password: rd.Password,
		DB:       rd.DB,
	})

	if err := rd.Client.Ping(context.Background()).Err(); err != nil {
		return err
	}

	logger.Info.Println("Connected to Redis successfully")

	return nil
}

// Reconnect attempts to reconnect to Redis using the retry utility.
func (rd *Client) Reconnect() {
	err := utils.Retry(5, 2*time.Second, rd.connect)
	if err != nil {
		logger.Error.Println("Failed to reconnect to Redis after multiple attempts")
	} else {
		logger.Info.Println("Reconnected to Redis successfully")
	}
}

// Close gracefully closes the Redis client connection.
func (rd *Client) Close() {
	if err := rd.Client.Close(); err != nil {
		logger.Warn.Println("Error closing Redis connection:", err)
	} else {
		logger.Info.Println("Redis connection closed")
	}
}
