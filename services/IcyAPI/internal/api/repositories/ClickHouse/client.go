package clickhouse

import (
	logger "itsjaylen/IcyLogger"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

// ClickHouseClient wraps the GORM DB instance for ClickHouse.
type ClickHouseClient struct {
	DB *gorm.DB
}

// NewClickHouseClient initializes and returns a ClickHouse client.
func NewClickHouseClient(dsn string) (*ClickHouseClient, error) {
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Ping the database to ensure connectivity
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	logger.Info.Println("Connected to ClickHouse successfully")
	return &ClickHouseClient{DB: db}, nil
}
