package routes

import (
	"IcyAPI/internal/api/routes/admin"
	"IcyAPI/internal/appinit"
	"net/http"
)

// Register all routes here
func InitRegisterRoutes(mux *http.ServeMux, app *appinit.App) {
	admin.RegisterRoutes(mux, app.RedisClient, app.PostgresClient)
}
