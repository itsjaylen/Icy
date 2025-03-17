package middleware

import (
	"IcyAPI/internal/api/middleware/auth"
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func RoleMiddleware(next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store claims in request context using the custom context key
		ctx := context.WithValue(r.Context(), auth.ClaimsContextKey, claims)
		r = r.WithContext(ctx)

		// Check if the role matches
		if slices.Contains(roles, claims.Role) {
			next(w, r)
			return
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}