package auth

import (
	"IcyAPI/internal/api/middleware"
	auth "IcyAPI/internal/api/middleware/auth"
	"IcyAPI/internal/appinit"
	"net/http"
	"time"
)

func RegisterRoutes(mux *http.ServeMux, app *appinit.App) {
	authhandler := auth.NewAuthService(app.RedisClient, app.PostgresClient)

	mux.HandleFunc("/signup", middleware.RateLimitMiddleware(authhandler.SignupHandler, 1*time.Second))
	mux.HandleFunc("/login", middleware.RateLimitMiddleware(authhandler.LoginHandler, 1*time.Second))
	mux.HandleFunc("/refresh", authhandler.RefreshTokenHandler)
	mux.HandleFunc("/admin", middleware.RoleMiddleware(auth.AdminHandler, "admin"))
	mux.HandleFunc("/user", middleware.RoleMiddleware(auth.UserHandler, "user", "admin"))
	mux.HandleFunc("/regen-api-key", middleware.RoleMiddleware(authhandler.RegenAPIKeyHandler, "user", "admin"))
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
