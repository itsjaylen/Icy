// Package appinit provides functions to initialize Redis dependencies.
package appinit

import (
	"fmt"

	redis "github.com/itsjaylen/IcyAPI/internal/api/repositories/Redis"
	config "itsjaylen/IcyConfig"
)

// InitRedis initializes a Redis client.
func InitRedis(cfg *config.AppConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	return redis.NewClient(addr, cfg.Redis.Password, 0)
}
