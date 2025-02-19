package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)


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

// Ping checks if Redis is responsive.
func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.Client.Ping(ctx).Result()
	return err
}

// TTL retrieves the time-to-live (TTL) of a key in Redis.
func (r *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

// MSet sets multiple key-value pairs in Redis.
func (r *RedisClient) MSet(ctx context.Context, keyValuePairs map[string]interface{}, expiration time.Duration) error {
	pipe := r.Client.Pipeline()
	for key, value := range keyValuePairs {
		pipe.Set(ctx, key, value, expiration)
	}
	_, err := pipe.Exec(ctx)
	return err
}

// MGet retrieves multiple values from Redis by key.
func (r *RedisClient) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return r.Client.MGet(ctx, keys...).Result()
}

// MDel removes multiple keys from Redis.
func (r *RedisClient) MDel(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Latency measures the response time of a Redis PING command.
func (r *RedisClient) Latency(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	_, err := r.Client.Ping(ctx).Result()
	if err != nil {
		return 0, err
	}
	return time.Since(start), nil
}
