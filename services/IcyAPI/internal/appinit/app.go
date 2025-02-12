package appinit

import (
	"fmt"

	clickhouse "IcyAPI/internal/api/repositories/ClickHouse"
	minobucket "IcyAPI/internal/api/repositories/MinoBucket"
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	rabbitmq "IcyAPI/internal/api/repositories/RabbitMQ"
	redis "IcyAPI/internal/api/repositories/Redis"
	"IcyAPI/internal/api/services/webhooks"
	"IcyAPI/internal/events"

	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"
)

// App structure to hold dependencies
type App struct {
	Cfg              *config.AppConfig
	EventServer      *events.EventServer
	RedisClient      *redis.RedisClient
	ClickHouseClient *clickhouse.ClickHouseClient
	PostgresClient   *postgresql.PostgresClient
	MinioClient      *minobucket.MinioClient
	RabbitMQ         *rabbitmq.RabbitMQClient
}

// NewApp initializes the application
func NewApp(debug bool) (*App, error) {
	logger.Debug.Printf("Debug mode: %v", debug)
	cfg, err := config.LoadConfig(map[bool]string{true: "debug", false: "release"}[debug])
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	if cfg.Webhook.Enabled {
		if err := webhooks.SendDiscordWebhook(cfg.Webhook.URL, fmt.Sprintf("IcyAPI has started! 🚀 in mode: %v", debug)); err != nil {
			logger.Error.Printf("Error sending Discord webhook: %v", err)
		}
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

	rabbitmqClient, err := InitRabbitMQ(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		Cfg:              cfg,
		EventServer:      events.NewEventServer(cfg.Server.Host, cfg.EventServer.Port),
		RedisClient:      redisClient,
		ClickHouseClient: clickhouseClient,
		PostgresClient:   postgresClient,
		MinioClient:      minioClient,
		RabbitMQ:         rabbitmqClient,
	}, nil
}
