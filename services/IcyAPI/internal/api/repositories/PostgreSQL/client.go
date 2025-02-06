package postgresql

import (
	logger "itsjaylen/IcyLogger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// PostgresClient wraps the GORM DB instance for PostgreSQL.
type PostgresClient struct {
	DB *gorm.DB
}

// NewPostgresClient initializes and returns a PostgreSQL client.
func NewPostgresClient(dsn string) (*PostgresClient, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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

	logger.Info.Println("Connected to PostgreSQL successfully")
	return &PostgresClient{DB: db}, nil
}
