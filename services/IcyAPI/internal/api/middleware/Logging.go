package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs HTTP requests similar to Gin's logger.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		ww := &responseWriterWrapper{ResponseWriter: writer, statusCode: http.StatusOK}
		next.ServeHTTP(ww, request)

		// Log request details
		duration := time.Since(start)
		log.Printf("%s %s %d %s - %s",
			request.Method,
			request.URL.Path,
			ww.statusCode,
			duration,
			request.UserAgent(),
		)
	})
}
