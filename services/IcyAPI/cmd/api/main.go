package main

import (
	"IcyAPI/internal/api/server"
	"IcyAPI/internal/events"
	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
)

// App structure to hold dependencies for both servers
type App struct {
	cfg         *config.AppConfig
	apiServer   *server.Server
	eventServer *events.EventServer
}

// NewApp initializes the application, both API and Event servers
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

	// Initialize both the API server and the Event server
	apiServer := server.NewAPIServer(cfg.Server.Host, cfg.Server.Port)
	eventServer := events.NewEventServer(cfg.Server.Host, cfg.EventServer.Port)

	return &App{
		cfg:         cfg,
		apiServer:   apiServer,
		eventServer: eventServer,
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

	// Start the API server in a goroutine
	go func() {
		if err := app.apiServer.Start(); err != nil {
			logger.Error.Printf("Error starting API server: %v", err)
			stop()
		}
	}()

	// Start the Event server in a goroutine
	go func() {
		if err := app.eventServer.Start(); err != nil {
			logger.Error.Printf("Error starting Event server: %v", err)
			stop()
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info.Println("Shutting down servers...")

	// Gracefully shutdown the API and Event servers
	app.apiServer.Shutdown()
	app.eventServer.Shutdown()

	logger.Info.Println("Servers gracefully stopped.")
}
