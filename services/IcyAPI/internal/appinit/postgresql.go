// Package appinit provides functions to initialize PostgreSQL dependencies.
package appinit

import (
	"fmt"

	postgresql "github.com/itsjaylen/IcyAPI/internal/api/repositories/PostgreSQL"
	config "itsjaylen/IcyConfig"
)

// InitPostgreSQL initializes a PostgreSQL client.
func InitPostgreSQL(cfg *config.AppConfig) (*postgresql.PostgresClient, error) {
	addr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	return postgresql.NewPostgresClient(addr)
}
