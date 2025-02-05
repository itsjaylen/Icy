package main

import (
	clickhouse "IcyAPI/internal/api/repositories/ClickHouse"
	minobucket "IcyAPI/internal/api/repositories/MinoBucket"
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"
	"IcyAPI/internal/api/server"
	"IcyAPI/internal/events"
	"context"
	"fmt"
	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
)

// App structure to hold dependencies for both servers
type App struct {
	cfg              *config.AppConfig
	apiServer        *server.Server
	eventServer      *events.EventServer
	redisClient      *redis.RedisClient
	clickhouseClient *clickhouse.ClickHouseClient
	postgresClient   *postgresql.PostgresClient
	miniobucketClient *minobucket.MinioClient
}

// NewApp initializes the application
func NewApp(debug bool) (*App, error) {
	var err error
	var cfg *config.AppConfig

	if debug {
		logger.Debug.Printf("Debug mode is enabled")
		cfg, err = config.LoadConfig("debug")
	} else {
		logger.Debug.Printf("Debug mode is disabled")
		cfg, err = config.LoadConfig("release")
	}

	if err != nil {
		logger.Error.Fatalf("Error loading config: %v", err)
	}

	// Initialize Redis
	redisAddress := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	redisClient, err := redis.NewRedisClient(redisAddress, cfg.Redis.Password, 0)
	if err != nil {
		logger.Error.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Initialize ClickHouse
	clickhouseAddress := fmt.Sprintf("clickhouse://%s:%s?username=%s&password=%s",
		cfg.Clickhouse.Host,
		cfg.Clickhouse.Port,
		cfg.Clickhouse.User,
		cfg.Clickhouse.Password,
	)

	clickhouseClient, err := clickhouse.NewClickHouseClient(clickhouseAddress)
	if err != nil {
		logger.Error.Fatalf("Failed to initialize ClickHouse: %v", err)
	}

	// Initialize PostgreSQL
	postgresAddress := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	postgresClient, err := postgresql.NewPostgresClient(postgresAddress)
	if err != nil {
		logger.Error.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	// Initialize Minio
	minioAddress := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)
	minioClient, err := minobucket.NewMinioClient(minioAddress, cfg.Minio.AccessKey, cfg.Minio.SecretKey, false)
	if err != nil {
		logger.Error.Fatalf("Failed to initialize Minio: %v", err)
	}

	// Initialize servers
	apiServer := server.NewAPIServer(cfg.Server.Host, cfg.Server.Port)
	eventServer := events.NewEventServer(cfg.Server.Host, cfg.EventServer.Port)

	return &App{
		cfg:              cfg,
		apiServer:        apiServer,
		eventServer:      eventServer,
		redisClient:      redisClient,
		clickhouseClient: clickhouseClient,
		postgresClient:   postgresClient,
		miniobucketClient: minioClient,
	}, nil
}

func main() {
	debug := pflag.Bool("debug", false, "Enable debug mode")
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.Parse()

	app, err := NewApp(*debug)
	if err != nil {
		logger.Error.Fatalf("Error initializing app: %v", err)
	}

	// Set up graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.apiServer.Start(); err != nil {
			logger.Error.Printf("Error starting API server: %v", err)
			stop()
		}
	}()

	go func() {
		if err := app.eventServer.Start(); err != nil {
			logger.Error.Printf("Error starting Event server: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info.Println("Shutting down servers...")

	// Graceful shutdown
	app.apiServer.Shutdown()
	app.eventServer.Shutdown()
	app.redisClient.Close()

	logger.Info.Println("Servers gracefully stopped.")
}
