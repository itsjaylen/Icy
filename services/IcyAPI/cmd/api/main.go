package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	appInit "IcyAPI/internal/init"
	logger "itsjaylen/IcyLogger"

	"github.com/spf13/pflag"
)

func main() {
	debug := pflag.Bool("debug", false, "Enable debug mode")
	pflag.Parse()

	app, err := appInit.NewApp(*debug)
	if err != nil {
		logger.Error.Fatalf("Failed to initialize app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go startServer("API", app.APIServer.Start, stop)
	go startServer("Event", app.EventServer.Start, stop)

	<-ctx.Done()
	logger.Info.Println("Shutting down servers...")

	if err := app.APIServer.Shutdown(); err != nil {
		logger.Error.Printf("Error shutting down API server: %v", err)
	}

	if err := app.EventServer.Shutdown(); err != nil {
		logger.Error.Printf("Error shutting down Event server: %v", err)
	}

	app.RedisClient.Close()
	logger.Info.Println("Servers gracefully stopped.")
}

func startServer(name string, startFunc func() error, stop context.CancelFunc) {
	if err := startFunc(); err != nil {
		logger.Error.Printf("Error starting %s server: %v", name, err)
		stop()
	}
}
