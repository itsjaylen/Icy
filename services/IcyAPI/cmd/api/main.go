// Package main is the entry point for the IcyAPI.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsjaylen/IcyAPI/internal/api/server"
	appInit "github.com/itsjaylen/IcyAPI/internal/appinit"
	"github.com/itsjaylen/IcyAPI/internal/workers"
	"github.com/spf13/pflag"
	logger "itsjaylen/IcyLogger"
)

func main() {
	debug := pflag.Bool("debug", false, "Enable debug mode")
	pflag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := appInit.NewApp(*debug)
	if err != nil {
		logger.Error.Fatalf("Failed to initialize app: %v", err)
	}

	// Run database migrations
	if err := app.RunMigrations(); err != nil {
		logger.Error.Fatalf("Database migration failed: %v", err)
	}

	apiServer := server.NewAPIServer(app)

	workers.SetupTaskManager(app)

	go startServer("API", apiServer.Start, stop)
	go startServer("Event", app.EventServer.Start, stop)

	<-ctx.Done()
	logger.Info.Println("Shutting down servers...")

	if err := apiServer.Shutdown(); err != nil {
		logger.Error.Printf("Error shutting down API server: %v", err)
	}

	if err := app.EventServer.Shutdown(); err != nil {
		logger.Error.Printf("Error shutting down Event server: %v", err)
	}

	app.Client.Close()
	logger.Info.Println("Servers gracefully stopped.")
}

func startServer(name string, startFunc func() error, stop context.CancelFunc) {
	if err := startFunc(); err != nil {
		logger.Error.Printf("Error starting %s server: %v", name, err)
		stop()
	}
}
