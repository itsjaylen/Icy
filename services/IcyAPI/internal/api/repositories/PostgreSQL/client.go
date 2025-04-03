// Package postgresql provides a PostgreSQL client with retry logic.
package postgresql

import (
	"time"

	"github.com/itsjaylen/IcyAPI/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	logger "itsjaylen/IcyLogger"
)

// PostgresClient wraps the GORM DB instance for PostgreSQL.
type PostgresClient struct {
	DB  *gorm.DB
	DSN string
}

// NewPostgresClient initializes and returns a PostgreSQL client with retry logic.
func NewPostgresClient(dsn string) (*PostgresClient, error) {
	client := &PostgresClient{DSN: dsn}

	err := utils.Retry(5, 2*time.Second, client.connect)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// connect establishes a connection to PostgreSQL and verifies connectivity.
func (pg *PostgresClient) connect() error {
	db, err := gorm.Open(postgres.Open(pg.DSN), &gorm.Config{})
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

	logger.Info.Println("Connected to PostgreSQL successfully")
	pg.DB = db

	return nil
}

// Reconnect attempts to reconnect to PostgreSQL using the retry utility.
func (pg *PostgresClient) Reconnect() {
	err := utils.Retry(5, 2*time.Second, pg.connect)
	if err != nil {
		logger.Error.Println("Failed to reconnect to PostgreSQL after multiple attempts")
	} else {
		logger.Info.Println("Reconnected to PostgreSQL successfully")
	}
}

// Migrate runs the database migrations for the provided models.
func (pg *PostgresClient) Migrate(models ...any) error {
	if err := pg.DB.AutoMigrate(models...); err != nil {
		return err
	}
	logger.Info.Println("Database migrations completed successfully")

	return nil
}
