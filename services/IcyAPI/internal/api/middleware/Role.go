package middleware

import (
	"IcyAPI/internal/api/middleware/auth"
	"IcyAPI/internal/appinit"
	"context"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func RoleMiddleware(config *appinit.App, next http.HandlerFunc, roles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")

		claims := &auth.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.Cfg.Server.JwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store claims in request context
		ctx := context.WithValue(r.Context(), auth.ClaimsContextKey, claims)
		r = r.WithContext(ctx)

		if slices.Contains(roles, claims.Role) {
			next(w, r)
			return
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}
