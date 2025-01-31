package middleware

import (
	"net/http"
	"runtime/debug"

	logger "itsjaylen/IcyLogger"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware handles panics 
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal Server Error",
					"message": "Something went wrong. Please try again later.",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}