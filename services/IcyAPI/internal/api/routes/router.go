package routes

import (
	"IcyAPI/internal/api/routes/admin"

	"github.com/gin-gonic/gin"
)

// Register all routes here
func InitRegisterRoutes(router *gin.Engine) {
	admin.RegisterRoutes(router)
}
