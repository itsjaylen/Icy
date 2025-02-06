package init

import (
	rabbitmq "IcyAPI/internal/api/repositories/RabbitMQ"
	"fmt"
	config "itsjaylen/IcyConfig"
)

func InitRabbitMQ(cfg *config.AppConfig) (*rabbitmq.RabbitMQClient, error) {
	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	return rabbitmq.NewRabbitMQClient(dsn)
}
