package routes

import (
	"IcyAPI/internal/api/routes/admin"
	"net/http"
)

// Register all routes here
func InitRegisterRoutes(mux *http.ServeMux) {
	admin.RegisterRoutes(mux)
}
