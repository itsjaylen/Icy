package rabbitmq

import (
	logger "itsjaylen/IcyLogger"

	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQClient(dsn string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	logger.Info.Println("Connected to RabbitMQ successfully")
	return &RabbitMQClient{Conn: conn, Channel: ch}, nil
}

func (r *RabbitMQClient) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
	logger.Info.Println("RabbitMQ connection closed")
}
