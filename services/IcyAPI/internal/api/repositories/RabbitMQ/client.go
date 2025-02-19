package rabbitmq

import (
	"IcyAPI/internal/utils"
	"time"

	logger "itsjaylen/IcyLogger"

	"github.com/streadway/amqp"
)

// RabbitMQClient wraps the RabbitMQ connection and channel.
type RabbitMQClient struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	DSN     string
}

// NewRabbitMQClient initializes and returns a RabbitMQ client with retry logic.
func NewRabbitMQClient(dsn string) (*RabbitMQClient, error) {
	client := &RabbitMQClient{DSN: dsn}

	err := utils.Retry(5, 2*time.Second, client.connect)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// connect establishes a connection to RabbitMQ and creates a channel.
func (r *RabbitMQClient) connect() error {
	conn, err := amqp.Dial(r.DSN)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	logger.Info.Println("Connected to RabbitMQ successfully")
	r.Conn = conn
	r.Channel = ch
	return nil
}

// Reconnect attempts to reconnect to RabbitMQ using the retry utility.
func (r *RabbitMQClient) Reconnect() {
	err := utils.Retry(5, 2*time.Second, r.connect)
	if err != nil {
		logger.Error.Println("Failed to reconnect to RabbitMQ after multiple attempts")
	} else {
		logger.Info.Println("Reconnected to RabbitMQ successfully")
	}
}

// Close gracefully closes the RabbitMQ connection and channel.
func (r *RabbitMQClient) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
	logger.Info.Println("RabbitMQ connection closed")
}
