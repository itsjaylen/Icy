package appinit

import (
	"fmt"

	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"

	clickhouse "github.com/itsjaylen/IcyAPI/internal/api/repositories/ClickHouse"
	minobucket "github.com/itsjaylen/IcyAPI/internal/api/repositories/MinoBucket"
	postgresql "github.com/itsjaylen/IcyAPI/internal/api/repositories/PostgreSQL"
	rabbitmq "github.com/itsjaylen/IcyAPI/internal/api/repositories/RabbitMQ"
	redis "github.com/itsjaylen/IcyAPI/internal/api/repositories/Redis"
	"github.com/itsjaylen/IcyAPI/internal/api/services/urlshortern"
	"github.com/itsjaylen/IcyAPI/internal/api/services/webhooks"
	"github.com/itsjaylen/IcyAPI/internal/events"
	"github.com/itsjaylen/IcyAPI/internal/models"
)

// App structure to hold dependencies.
type App struct {
	Cfg              *config.AppConfig
	EventServer      *events.EventServer
	Client           *redis.Client
	ClickHouseClient *clickhouse.Client
	PostgresClient   *postgresql.PostgresClient
	MinioClient      *minobucket.MinioClient
	RabbitMQ         *rabbitmq.Client
}

// NewApp initializes the application.
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

	Client, err := InitRedis(cfg)
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
		Client:           Client,
		ClickHouseClient: clickhouseClient,
		PostgresClient:   postgresClient,
		MinioClient:      minioClient,
		RabbitMQ:         rabbitmqClient,
	}, nil
}

// RunMigrations handles database migrations.
func (a *App) RunMigrations() error {
	return a.PostgresClient.Migrate(&models.User{}, urlshortern.URLMapping{})
}
