// Package middleware ErrorHandler is a middleware for handling errors.
package middleware

import (
	"net/http"
	"runtime/debug"

	logger "itsjaylen/IcyLogger"
)

// RecoveryMiddleware handles panics.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error.Printf("Error: %v\n%s", err, debug.Stack())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
