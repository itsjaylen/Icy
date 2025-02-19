package clickhouse

import (
	"IcyAPI/internal/utils"
	"time"

	"itsjaylen/IcyLogger"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

// ClickHouseClient wraps the GORM DB instance for ClickHouse.
type ClickHouseClient struct {
	DB  *gorm.DB
	DSN string
}

// NewClickHouseClient initializes and returns a ClickHouse client with retry logic.
func NewClickHouseClient(dsn string) (*ClickHouseClient, error) {
	client := &ClickHouseClient{DSN: dsn}
	err := utils.Retry(5, 2*time.Second, client.connect)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// connect establishes a connection to ClickHouse.
func (c *ClickHouseClient) connect() error {
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
func (c *ClickHouseClient) Reconnect() {
	err := utils.Retry(5, 2*time.Second, c.connect)
	if err != nil {
		logger.Error.Println("Failed to reconnect to ClickHouse after multiple attempts")
	} else {
		logger.Info.Println("Reconnected to ClickHouse successfully")
	}
}