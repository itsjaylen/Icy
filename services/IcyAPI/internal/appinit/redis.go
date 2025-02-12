package appinit

import (
	redis "IcyAPI/internal/api/repositories/Redis"
	"fmt"
	config "itsjaylen/IcyConfig"
)

func InitRedis(cfg *config.AppConfig) (*redis.RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	return redis.NewRedisClient(addr, cfg.Redis.Password, 0)
}
