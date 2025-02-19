package server

import (
	"IcyAPI/internal/api/middleware"
	"IcyAPI/internal/api/routes"
	"IcyAPI/internal/appinit"
	"IcyAPI/internal/workers/tasks/health"
	"context"
	"fmt"
	logger "itsjaylen/IcyLogger"
	"net/http"
	"time"
)

// Server struct to hold server configurations
type Server struct {
	Host    string
	Port    string
	Handler http.Handler
	server  *http.Server
}

// NewAPIServer creates a new server instance with injected dependencies
func NewAPIServer(app *appinit.App) *Server {
	mux := http.NewServeMux()

	// Register healthz endpoint
	mux.HandleFunc("/healthz", health.HealthzHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register routes with dependencies
	routes.InitRegisterRoutes(mux, app)

	// Apply middlewares
	handler := middleware.LoggingMiddleware(mux)
	handler = middleware.AnalyticsMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)
	handler = middleware.ErrorHandler(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", app.Cfg.Server.Host, app.Cfg.Server.Port),
		Handler: handler,
	}

	return &Server{
		Host:    app.Cfg.Server.Host,
		Port:    app.Cfg.Server.Port,
		Handler: handler,
		server:  srv,
	}
}


// Start runs the server
func (s *Server) Start() error {
	logger.Info.Printf("Starting server on %s:%s", s.Host, s.Port)
	return s.server.ListenAndServe()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
