package main

import (
	"IcyAPI/internal/api/server"
	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
)

// App structure to hold dependencies
 type App struct {
	cfg    *config.AppConfig
	server *server.Server
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

	server := server.NewServer(cfg.Server.Host, cfg.Server.Port)

	return &App{
		cfg:    cfg,
		server: server,
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

	logger.Info.Printf("Loaded config: %v", app.cfg)

	// Create a context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run the server in a goroutine
	go func() {
		if err := app.server.Start(); err != nil {
			logger.Error.Printf("Error starting server: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info.Println("Shutting down server...")

	if err := app.server.Shutdown(); err != nil {
		logger.Error.Printf("Error during server shutdown: %v", err)
	}

	logger.Info.Println("Server gracefully stopped.")
}