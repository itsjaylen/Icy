package middleware

import (
	"net/http"

	logger "itsjaylen/IcyLogger"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if addErr := c.Error(err.Err); addErr != nil {
					logger.Error.Printf("Failed to add error to context: %v", addErr)
				}
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": c.Errors[0].Error(),
			})

			return
		}
	}
}