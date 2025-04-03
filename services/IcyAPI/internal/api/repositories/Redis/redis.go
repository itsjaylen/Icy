package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// Set stores a key-value pair in Redis with an optional expiration time.
func (rd *Client) Set(
	ctx context.Context,
	key string,
	value string,
	expiration time.Duration,
) error {
	return rd.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis by key.
func (rd *Client) Get(ctx context.Context, key string) (string, error) {
	value, err := rd.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}

	return value, err
}

// Exists checks if a key exists in Redis.
func (rd *Client) Exists(ctx context.Context, key string) (bool, error) {
	count, err := rd.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Delete removes a key from Redis.
func (rd *Client) Delete(ctx context.Context, key string) error {
	return rd.Client.Del(ctx, key).Err()
}

// SetJSON stores a struct in Redis as a JSON string.
func (rd *Client) SetJSON(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rd.Set(ctx, key, string(jsonData), expiration)
}

// GetJSON retrieves a JSON-stored value and unmarshals it into a struct.
func (rd *Client) GetJSON(ctx context.Context, key string, dest interface{}) error {
	jsonStr, err := rd.Get(ctx, key)
	if err != nil {
		return err
	}
	if jsonStr == "" {
		return nil
	}

	return json.Unmarshal([]byte(jsonStr), dest)
}

// Ping checks if Redis is responsive.
func (rd *Client) Ping(ctx context.Context) error {
	_, err := rd.Client.Ping(ctx).Result()

	return err
}

// TTL retrieves the time-to-live (TTL) of a key in Redis.
func (rd *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return rd.Client.TTL(ctx, key).Result()
}

// MSet sets multiple key-value pairs in Redis.
func (rd *Client) MSet(ctx context.Context, keyValuePairs map[string]interface{}, expiration time.Duration) error {
	pipe := rd.Client.Pipeline()
	for key, value := range keyValuePairs {
		pipe.Set(ctx, key, value, expiration)
	}
	_, err := pipe.Exec(ctx)

	return err
}

// MGet retrieves multiple values from Redis by key.
func (rd *Client) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return rd.Client.MGet(ctx, keys...).Result()
}

// MDel removes multiple keys from Redis.
func (rd *Client) MDel(ctx context.Context, keys ...string) error {
	return rd.Client.Del(ctx, keys...).Err()
}

// Latency measures the response time of a Redis PING command.
func (rd *Client) Latency(ctx context.Context) (time.Duration, error) {
	start := time.Now()
	_, err := rd.Client.Ping(ctx).Result()
	if err != nil {
		return 0, err
	}

	return time.Since(start), nil
}
