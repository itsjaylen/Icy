// Package auth contains the models used for authentication and authorization.
package auth

import (
	"github.com/golang-jwt/jwt/v5"
	postgresql "github.com/itsjaylen/IcyAPI/internal/api/repositories/PostgreSQL"
	redis "github.com/itsjaylen/IcyAPI/internal/api/repositories/Redis"
	config "itsjaylen/IcyConfig"
)

// UserSignupRequest represents the request payload for user signup.
type UserSignupRequest struct {
	Username string `json:"username" example:"newuser" binding:"required"`           // Username of the user
	Password string `json:"password" example:"securepassword123" binding:"required"` // Password for the user
}

// ErrorResponse represents a generic error response message.
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// SuccessResponse represents a generic success message.
type SuccessResponse struct {
	Message string `json:"message"`
}

// LoginResponse represents the response containing authentication tokens.
type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	APIKey       string `json:"apiKey"`
}

// Claims represents JWT claims including username and role.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Service handles authentication and user-related operations.
type Service struct {
	PostgresClient *postgresql.PostgresClient
	Client         *redis.Client
	Config         *config.AppConfig
}
