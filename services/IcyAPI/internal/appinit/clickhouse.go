package appinit

import (
	clickhouse "IcyAPI/internal/api/repositories/ClickHouse"
	"fmt"
	config "itsjaylen/IcyConfig"
)

func InitClickHouse(cfg *config.AppConfig) (*clickhouse.ClickHouseClient, error) {
	addr := fmt.Sprintf("clickhouse://%s:%s?username=%s&password=%s",
		cfg.Clickhouse.Host, cfg.Clickhouse.Port, cfg.Clickhouse.User, cfg.Clickhouse.Password)
	return clickhouse.NewClickHouseClient(addr)
}
