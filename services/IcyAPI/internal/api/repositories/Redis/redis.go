package redis

import (
	"context"
	"encoding/json"
	logger "itsjaylen/IcyLogger"
	"time"

	"github.com/redis/go-redis/v9"
)

// Close shuts down the Redis connection.
func (r *RedisClient) Close() {
	if err := r.Client.Close(); err != nil {
		logger.Error.Printf("Error closing Redis connection: %v", err)
	} else {
		logger.Info.Println("Redis connection closed")
	}
}

// Set stores a key-value pair in Redis with an optional expiration time.
func (r *RedisClient) Set(
	ctx context.Context,
	key string,
	value string,
	expiration time.Duration,
) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key.
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return value, err
}

// Exists checks if a key exists in Redis.
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Delete removes a key from Redis.
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// SetJSON stores a struct in Redis as a JSON string.
func (r *RedisClient) SetJSON(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Set(ctx, key, string(jsonData), expiration)
}

// GetJSON retrieves a JSON-stored value and unmarshals it into a struct.
func (r *RedisClient) GetJSON(ctx context.Context, key string, dest interface{}) error {
	jsonStr, err := r.Get(ctx, key)
	if err != nil {
		return err
	}
	if jsonStr == "" {
		return nil
	}
	return json.Unmarshal([]byte(jsonStr), dest)
}
