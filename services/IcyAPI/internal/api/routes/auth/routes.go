// Package auth provides routes for user authentication and authorization.
package auth

import (
	"net/http"
	"time"

	"github.com/itsjaylen/IcyAPI/internal/api/middleware"
	auth "github.com/itsjaylen/IcyAPI/internal/api/middleware/auth"
	"github.com/itsjaylen/IcyAPI/internal/appinit"
)

// RegisterRoutes registers routes for user authentication and authorization. TODO: Complete the ouath.
func RegisterRoutes(mux *http.ServeMux, app *appinit.App) {
	authhandler := auth.NewAuthService(app.Client, app.PostgresClient, app.Cfg)

	mux.HandleFunc("/signup", middleware.RateLimitMiddleware(authhandler.SignupHandler, 5*time.Second, 3))
	mux.HandleFunc("/login", middleware.RateLimitMiddleware(authhandler.LoginHandler, 5*time.Second, 3))
	mux.HandleFunc("/refresh", authhandler.RefreshTokenHandler)
	mux.HandleFunc("/admin", middleware.RoleMiddleware(app, auth.AdminHandler, "admin"))
	mux.HandleFunc("/user", middleware.RoleMiddleware(app, auth.UserHandler, "user", "admin"))
	mux.HandleFunc("/regen-api-key", middleware.RoleMiddleware(app, authhandler.RegenAPIKeyHandler, "user", "admin"))

	mux.HandleFunc("/logout", authhandler.LogoutHandler)

	// Google OAuth
	mux.HandleFunc("/login/google", auth.GoogleLoginHandler)
	mux.HandleFunc("/callback/google", authhandler.GoogleCallbackHandler)

	// GitHub OAuth
	mux.HandleFunc("/login/github", auth.GithubLoginHandler)
	mux.HandleFunc("/callback/github", authhandler.GithubCallbackHandler)

	// Twitch OAuth
	mux.HandleFunc("/login/twitch", auth.TwitchLoginHandler)
	mux.HandleFunc("/callback/twitch", authhandler.TwitchCallbackHandler)
}
