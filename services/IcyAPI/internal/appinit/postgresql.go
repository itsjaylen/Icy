package appinit

import (
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	"fmt"
	config "itsjaylen/IcyConfig"
)

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
