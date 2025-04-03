// Package admin provides routes for administrative tasks.
package admin

import (
	"net/http"
	"time"

	"github.com/itsjaylen/IcyAPI/internal/api/controllers/admin"
	"github.com/itsjaylen/IcyAPI/internal/api/middleware"
	postgresql "github.com/itsjaylen/IcyAPI/internal/api/repositories/PostgreSQL"
	redis "github.com/itsjaylen/IcyAPI/internal/api/repositories/Redis"
)

// RegisterRoutes registers admin-related routes. TODO: ADD AUTH!!!!!!
func RegisterRoutes(mux *http.ServeMux, client *redis.Client, postgresClient *postgresql.PostgresClient) {
	AdminController := admin.NewAdminController(client, postgresClient)

	mux.Handle(
		"/admin/users",
		middleware.RateLimiter(http.HandlerFunc(AdminController.HandleUserRequest), 5, 10*time.Second),
	)
	mux.Handle("/admin/status", http.HandlerFunc(admin.StatusHandler))
	mux.Handle("/admin/restart", http.HandlerFunc(admin.RestartHandler))
	mux.Handle("/admin/exec", http.HandlerFunc(admin.ExecHandler))
}
