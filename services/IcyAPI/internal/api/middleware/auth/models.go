package auth

import (
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"

	"github.com/golang-jwt/jwt/v5"
)

type UserSignupRequest struct {
	Username string `json:"username" example:"newuser" binding:"required"`           // Username of the user
	Password string `json:"password" example:"securepassword123" binding:"required"` // Password for the user
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	APIKey       string `json:"api_key"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthService struct {
	PostgresClient *postgresql.PostgresClient
	RedisClient    *redis.RedisClient
}
