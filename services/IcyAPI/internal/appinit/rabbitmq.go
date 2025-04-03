// Package appinit provides functions to initialize RabbitMQ dependencies.
package appinit

import (
	"fmt"

	rabbitmq "github.com/itsjaylen/IcyAPI/internal/api/repositories/RabbitMQ"
	config "itsjaylen/IcyConfig"
)

// InitRabbitMQ initializes a RabbitMQ client.
func InitRabbitMQ(cfg *config.AppConfig) (*rabbitmq.Client, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	return rabbitmq.NewClient(dsn)
}
