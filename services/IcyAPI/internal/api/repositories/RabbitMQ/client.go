// Package rabbitmq provides a RabbitMQ client with retry logic.
package rabbitmq

import (
	"time"

	"github.com/itsjaylen/IcyAPI/internal/utils"
	"github.com/streadway/amqp"
	logger "itsjaylen/IcyLogger"
)

// Client wraps the RabbitMQ connection and channel.
type Client struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	DSN     string
}

// NewClient initializes and returns a RabbitMQ client with retry logic.
func NewClient(dsn string) (*Client, error) {
	client := &Client{DSN: dsn}

	err := utils.Retry(5, 2*time.Second, client.connect)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// connect establishes a connection to RabbitMQ and creates a channel.
func (rabbitmq *Client) connect() error {
	conn, err := amqp.Dial(rabbitmq.DSN)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()

		return err
	}

	logger.Info.Println("Connected to RabbitMQ successfully")
	rabbitmq.Conn = conn
	rabbitmq.Channel = ch

	return nil
}

// Close gracefully closes the RabbitMQ connection and channel.
func (rabbitmq *Client) Close() {
	if rabbitmq.Channel != nil {
		if err := rabbitmq.Channel.Close(); err != nil {
			logger.Error.Println("Failed to close RabbitMQ channel:", err)
		}
	}
	if rabbitmq.Conn != nil {
		if err := rabbitmq.Conn.Close(); err != nil {
			logger.Error.Println("Failed to close RabbitMQ connection:", err)
		}
	}
	logger.Info.Println("RabbitMQ connection closed")
}

// Reconnect attempts to reconnect to RabbitMQ using the retry utility.
func (rabbitmq *Client) Reconnect() {
	err := utils.Retry(5, 2*time.Second, rabbitmq.connect)
	if err != nil {
		logger.Error.Println("Failed to reconnect to RabbitMQ after multiple attempts")
	} else {
		logger.Info.Println("Reconnected to RabbitMQ successfully")
	}
}
