package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	var requestCounts = make(map[string]int)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		mu.Lock()
		defer mu.Unlock()

		// Increment request count
		requestCounts[clientIP]++

		// Reset the counter after the window
		go func() {
			time.Sleep(window)
			mu.Lock()
			delete(requestCounts, clientIP)
			mu.Unlock()
		}()

		if requestCounts[clientIP] > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}