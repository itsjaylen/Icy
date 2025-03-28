package admin

import (
	"IcyAPI/internal/api/controllers/admin"
	"IcyAPI/internal/api/middleware"
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"
	"net/http"
	"time"
)

// RegisterRoutes registers admin-related routes
func RegisterRoutes(mux *http.ServeMux, redisClient *redis.RedisClient, postgresClient *postgresql.PostgresClient) {
	AdminController := admin.NewAdminController(redisClient, postgresClient)


	mux.Handle(
		"/admin/users",
		middleware.RateLimiter(http.HandlerFunc(AdminController.HandleUserRequest), 5, 10*time.Second),
	)
	mux.Handle("/admin/status", http.HandlerFunc(admin.StatusHandler))
	mux.Handle("/admin/restart", http.HandlerFunc(admin.RestartHandler))
	mux.Handle("/admin/exec", http.HandlerFunc(admin.ExecHandler))
}
