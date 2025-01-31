package server

import (
	"IcyAPI/internal/api/middleware"
	"IcyAPI/internal/api/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

// Server struct to hold server-related configurations
type Server struct {
	Host   string
	Port   string
	Router *gin.Engine
}

// NewServer creates a new server instance with the given host and port
func NewServer(host, port string) *Server {
	router := gin.Default()

	// Apply rate limiter middleware globally
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.ErrorHandler())

	// Register admin routes from the router package
	routes.InitRegisterRoutes(router)

	return &Server{
		Host:   host,
		Port:   port,
		Router: router,
	}
}

// Start runs the server on the given port
func (s *Server) Start() error {
	address := fmt.Sprintf("%s:%s", s.Host, s.Port)
	return s.Router.Run(address)
}
