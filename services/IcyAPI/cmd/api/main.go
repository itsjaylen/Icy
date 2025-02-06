package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	clickhouse "IcyAPI/internal/api/repositories/ClickHouse"
	minobucket "IcyAPI/internal/api/repositories/MinoBucket"
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"
	"IcyAPI/internal/api/server"
	"IcyAPI/internal/events"

	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"

	"github.com/spf13/pflag"
)

// App structure to hold dependencies
type App struct {
	cfg              *config.AppConfig
	apiServer        *server.Server
	eventServer      *events.EventServer
	redisClient      *redis.RedisClient
	clickhouseClient *clickhouse.ClickHouseClient
	postgresClient   *postgresql.PostgresClient
	minioClient      *minobucket.MinioClient
}

// NewApp initializes the application
func NewApp(debug bool) (*App, error) {
	logger.Debug.Printf("Debug mode: %v", debug)
	cfg, err := config.LoadConfig(map[bool]string{true: "debug", false: "release"}[debug])
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	redisClient, err := initRedis(cfg)
	if err != nil {
		return nil, err
	}

	clickhouseClient, err := initClickHouse(cfg)
	if err != nil {
		return nil, err
	}

	postgresClient, err := initPostgreSQL(cfg)
	if err != nil {
		return nil, err
	}

	minioClient, err := initMinio(cfg)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:              cfg,
		apiServer:        server.NewAPIServer(cfg.Server.Host, cfg.Server.Port),
		eventServer:      events.NewEventServer(cfg.Server.Host, cfg.EventServer.Port),
		redisClient:      redisClient,
		clickhouseClient: clickhouseClient,
		postgresClient:   postgresClient,
		minioClient:      minioClient,
	}, nil
}

func initRedis(cfg *config.AppConfig) (*redis.RedisClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	return redis.NewRedisClient(addr, cfg.Redis.Password, 0)
}

func initClickHouse(cfg *config.AppConfig) (*clickhouse.ClickHouseClient, error) {
	addr := fmt.Sprintf("clickhouse://%s:%s?username=%s&password=%s",
		cfg.Clickhouse.Host, cfg.Clickhouse.Port, cfg.Clickhouse.User, cfg.Clickhouse.Password)
	return clickhouse.NewClickHouseClient(addr)
}

func initPostgreSQL(cfg *config.AppConfig) (*postgresql.PostgresClient, error) {
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

func initMinio(cfg *config.AppConfig) (*minobucket.MinioClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)
	return minobucket.NewMinioClient(addr, cfg.Minio.AccessKey, cfg.Minio.SecretKey, false)
}

func main() {
	debug := pflag.Bool("debug", false, "Enable debug mode")
	pflag.Parse()

	app, err := NewApp(*debug)
	if err != nil {
		logger.Error.Fatalf("Failed to initialize app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go startServer("API", app.apiServer.Start, stop)
	go startServer("Event", app.eventServer.Start, stop)

	<-ctx.Done()
	logger.Info.Println("Shutting down servers...")

	if err := app.apiServer.Shutdown(); err != nil {
		logger.Error.Printf("Error shutting down API server: %v", err)
	}

	if err := app.eventServer.Shutdown(); err != nil {
		logger.Error.Printf("Error shutting down Event server: %v", err)
	}

	app.redisClient.Close()
	logger.Info.Println("Servers gracefully stopped.")
}

func startServer(name string, startFunc func() error, stop context.CancelFunc) {
	if err := startFunc(); err != nil {
		logger.Error.Printf("Error starting %s server: %v", name, err)
		stop()
	}
}
