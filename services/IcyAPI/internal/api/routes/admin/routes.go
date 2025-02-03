package admin

import (
	"IcyAPI/internal/api/controllers/admin"
	"IcyAPI/internal/api/middleware"
	"net/http"
	"time"
)

// RegisterRoutes registers admin-related routes
func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/admin/status", middleware.RateLimiter(http.HandlerFunc(admin.GetStatusHandler), 5, 10*time.Second))
}
