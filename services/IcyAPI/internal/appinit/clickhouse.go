// Package appinit provides functions to initialize ClickHouse dependencies.
package appinit

import (
	"fmt"

	clickhouse "github.com/itsjaylen/IcyAPI/internal/api/repositories/ClickHouse"
	config "itsjaylen/IcyConfig"
)

// InitClickHouse initializes a ClickHouse client.
func InitClickHouse(cfg *config.AppConfig) (*clickhouse.Client, error) {
	addr := fmt.Sprintf("clickhouse://%s:%s?username=%s&password=%s",
		cfg.Clickhouse.Host, cfg.Clickhouse.Port, cfg.Clickhouse.User, cfg.Clickhouse.Password)

	return clickhouse.NewClickHouseClient(addr)
}
