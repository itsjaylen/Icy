// Package middleware provides HTTP middleware functionalities such as authentication and role-based access control.
package middleware

import (
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsjaylen/IcyAPI/internal/api/middleware/auth"
	"github.com/itsjaylen/IcyAPI/internal/appinit"
)

// RoleMiddleware is an authentication middleware that checks the user's role before allowing access to a handler.
func RoleMiddleware(config *appinit.App, next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		tokenString := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")

		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
			return config.Cfg.Server.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)

			return
		}

		// Store claims in request context
		ctx := context.WithValue(request.Context(), auth.ClaimsContextKey, claims)
		request = request.WithContext(ctx)

		if slices.Contains(roles, claims.Role) {
			next(writer, request)

			return
		}

		http.Error(writer, "Forbidden", http.StatusForbidden)
	}
}
