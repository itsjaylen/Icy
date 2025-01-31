package admin

import "github.com/gin-gonic/gin"

func GetStatusHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}