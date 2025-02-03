package models

import "github.com/spf13/pflag"

// RabbitMQConfig holds the configuration for RabbitMQ connection.
type RabbitMQConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Enabled  bool   `json:"enabled"`
}

// RabbitMQFlags defines the command line flags for RabbitMQ configuration.
func RabbitMQFlags(fs *pflag.FlagSet) {
	fs.String("rabbitmq.host", "", "RabbitMQ host")
	fs.String("rabbitmq.port", "", "RabbitMQ port")
	fs.String("rabbitmq.user", "", "RabbitMQ user")
	fs.String("rabbitmq.password", "", "RabbitMQ password")
	fs.Bool("rabbitmq.enabled", false, "Enable RabbitMQ")
}
