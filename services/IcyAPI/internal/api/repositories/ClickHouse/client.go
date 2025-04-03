// Package clickhouse provides a ClickHouse client with retry logic.
package clickhouse

import (
	"time"

	"github.com/itsjaylen/IcyAPI/internal/utils"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	logger "itsjaylen/IcyLogger"
)

// Client represents a ClickHouse client.
type Client struct {
	DB  *gorm.DB
	DSN string
}

// NewClickHouseClient initializes and returns a ClickHouse client with retry logic.
func NewClickHouseClient(dsn string) (*Client, error) {
	client := &Client{DSN: dsn}
	err := utils.Retry(5, 2*time.Second, client.connect)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// connect establishes a connection to ClickHouse.
func (c *Client) connect() error {
	db, err := gorm.Open(clickhouse.Open(c.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	if err = sqlDB.Ping(); err != nil {
		return err
	}

	logger.Info.Println("Connected to ClickHouse successfully")
	c.DB = db

	return nil
}

// Reconnect attempts to reconnect to ClickHouse using the retry utility.
func (c *Client) Reconnect() {
	err := utils.Retry(5, 2*time.Second, c.connect)
	if err != nil {
		logger.Error.Println("Failed to reconnect to ClickHouse after multiple attempts")
	} else {
		logger.Info.Println("Reconnected to ClickHouse successfully")
	}
}
