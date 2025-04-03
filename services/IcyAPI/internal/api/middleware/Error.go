// Package middleware ErrorHandler is a middleware for handling errors.
package middleware

import (
	"net/http"
	"runtime/debug"

	logger "itsjaylen/IcyLogger"
)

// ErrorHandler is a middleware for handling errors.
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error.Printf("Error: %v\n%s", rec, debug.Stack())
				http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
