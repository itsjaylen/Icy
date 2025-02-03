package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests similar to Gin's logger
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		ww := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(ww, r)

		// Log request details
		duration := time.Since(start)
		log.Printf("%s %s %d %s - %s",
			r.Method,
			r.URL.Path,
			ww.statusCode,
			duration,
			r.UserAgent(),
		)
	})
}
