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
	userController := admin.NewUserController(redisClient, postgresClient)

	mux.Handle(
		"/admin/status",
		middleware.RateLimiter(http.HandlerFunc(admin.GetStatusHandler), 5, 10*time.Second),
	)
	mux.Handle(
		"/admin/users",
		middleware.RateLimiter(http.HandlerFunc(userController.HandleUserRequest), 5, 10*time.Second),
	)
}
