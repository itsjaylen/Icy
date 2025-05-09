// Package server provides the API server implementation.
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	logger "itsjaylen/IcyLogger"

	"github.com/itsjaylen/IcyAPI/internal/api/middleware"
	"github.com/itsjaylen/IcyAPI/internal/api/routes"
	"github.com/itsjaylen/IcyAPI/internal/appinit"
	"github.com/itsjaylen/IcyAPI/internal/workers/tasks/health"
	"github.com/rs/cors"
)

// Server struct to hold server configurations.
type Server struct {
	Handler http.Handler
	server  *http.Server
	Host    string
	Port    string
}

// NewAPIServer creates a new server instance with injected dependencies.
func NewAPIServer(app *appinit.App) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", health.HealthzHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			logger.Error.Printf("Failed to write response: %v", err)
		}
	})

	// Register routes with dependencies
	routes.InitRegisterRoutes(mux, app)

	// Apply middlewares
	handler := middleware.LoggingMiddleware(mux)
	// TODO: Fix this later: handler = middleware.AnalyticsMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)
	handler = middleware.ErrorHandler(handler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // or restrict to specific services consuming your API
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}).Handler(handler)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", app.Cfg.Server.Host, app.Cfg.Server.Port),
		Handler:           corsHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Server{
		Host:    app.Cfg.Server.Host,
		Port:    app.Cfg.Server.Port,
		Handler: corsHandler,
		server:  srv,
	}
}

// Start runs the server.
func (server *Server) Start() error {
	logger.Info.Printf("Starting server on %s:%s", server.Host, server.Port)

	return server.server.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (server *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.server.Shutdown(ctx)
}
