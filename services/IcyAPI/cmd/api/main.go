package main

import (
	"IcyAPI/internal/api/server"
	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"

	"github.com/spf13/pflag"
)

// App structure to hold dependencies such as config and server
type App struct {
	cfg    *config.AppConfig
	server *server.Server 
}

// Function to load the config and return the app structure
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

	// Initialize the server 
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

	// Create the app and inject the config
	app, err := NewApp(*debug)
	if err != nil {
		logger.Error.Fatalf("Error initializing app: %v", err)
	}

	// You can now access app.cfg anywhere in the app
	logger.Info.Printf("Loaded config: %v", app.cfg)

	if err := app.server.Start(); err != nil {
		logger.Error.Fatalf("Error starting server: %v", err)
	}
}
