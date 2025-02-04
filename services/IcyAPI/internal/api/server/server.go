package server

import (
	"IcyAPI/internal/api/middleware"
	"IcyAPI/internal/api/routes"
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

// NewServer creates a new server instance
func NewAPIServer(host, port string) *Server {
	mux := http.NewServeMux()

	// Register routes
	routes.InitRegisterRoutes(mux)

	// Apply middlewares
	handler := middleware.LoggingMiddleware(mux)
	handler = middleware.AnalyticsMiddleware(handler)
	handler = middleware.RecoveryMiddleware(handler)
	handler = middleware.ErrorHandler(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: handler,
	}

	return &Server{
		Host:    host,
		Port:    port,
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
