package init

import (
	"fmt"

	clickhouse "IcyAPI/internal/api/repositories/ClickHouse"
	minobucket "IcyAPI/internal/api/repositories/MinoBucket"
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"
	"IcyAPI/internal/api/server"
	"IcyAPI/internal/events"

	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"
)

// App structure to hold dependencies
type App struct {
	Cfg              *config.AppConfig
	APIServer        *server.Server
	EventServer      *events.EventServer
	RedisClient      *redis.RedisClient
	ClickHouseClient *clickhouse.ClickHouseClient
	PostgresClient   *postgresql.PostgresClient
	MinioClient      *minobucket.MinioClient
}

// NewApp initializes the application
func NewApp(debug bool) (*App, error) {
	logger.Debug.Printf("Debug mode: %v", debug)
	cfg, err := config.LoadConfig(map[bool]string{true: "debug", false: "release"}[debug])
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	redisClient, err := InitRedis(cfg)
	if err != nil {
		return nil, err
	}

	clickhouseClient, err := InitClickHouse(cfg)
	if err != nil {
		return nil, err
	}

	postgresClient, err := InitPostgreSQL(cfg)
	if err != nil {
		return nil, err
	}

	minioClient, err := InitMinio(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Cfg:              cfg,
		APIServer:        server.NewAPIServer(cfg.Server.Host, cfg.Server.Port),
		EventServer:      events.NewEventServer(cfg.Server.Host, cfg.EventServer.Port),
		RedisClient:      redisClient,
		ClickHouseClient: clickhouseClient,
		PostgresClient:   postgresClient,
		MinioClient:      minioClient,
	}, nil
}
