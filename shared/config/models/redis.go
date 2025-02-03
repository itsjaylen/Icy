package models

import "github.com/spf13/pflag"

// RedisConfig holds the configuration for Redis connection.
type RedisConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

// RedisFlags defines the command line flags for Redis configuration.
func RedisFlags(fs *pflag.FlagSet) {
	fs.String("redis.host", "", "Redis host")
	fs.String("redis.port", "", "Redis port")
	fs.String("redis.password", "", "Redis password")
}
