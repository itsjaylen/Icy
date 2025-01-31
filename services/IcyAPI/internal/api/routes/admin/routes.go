package admin

import (
	"IcyAPI/internal/api/middleware"
	"IcyAPI/internal/api/controllers/admin"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers user-related routes
func RegisterRoutes(router *gin.Engine) {
	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middleware.RateLimiter(5, time.Second*10))
	{
		adminRoutes.GET("/status", admin.GetStatusHandler)
	}
}