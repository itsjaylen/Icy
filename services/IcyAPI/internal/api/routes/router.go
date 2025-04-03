package routes

import (
	"net/http"

	"github.com/itsjaylen/IcyAPI/internal/api/routes/admin"
	"github.com/itsjaylen/IcyAPI/internal/api/routes/auth"
	"github.com/itsjaylen/IcyAPI/internal/appinit"
)

// Register all routes here.
func InitRegisterRoutes(mux *http.ServeMux, app *appinit.App) {
	// Register authentication routes
	auth.RegisterRoutes(mux, app)

	// Register admin routes
	admin.RegisterRoutes(mux, app.Client, app.PostgresClient)
}
