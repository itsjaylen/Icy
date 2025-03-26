package auth

import (
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"
	"IcyAPI/internal/models"
	"IcyAPI/internal/utils"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	config "itsjaylen/IcyConfig"
	logger "itsjaylen/IcyLogger"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/twitch"
	"gorm.io/gorm"
)

var (
	// JwtSecret          = []byte("SsYMx7hdNterwN011bzykWrMxjymmiu6")
	adminWhitelist     = map[string]bool{"superadmin": true}
	ClaimsContextKey   = contextKey("claims")
	oauth2GoogleConfig oauth2.Config
	oauth2GithubConfig oauth2.Config
	oauth2TwitchConfig oauth2.Config
	oauth2StateString  = "random"
)

func init() {
	// Google OAuth2 Config
	oauth2GoogleConfig = oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/callback/google",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	// GitHub OAuth2 Config
	oauth2GithubConfig = oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/callback/github",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// Twitch OAuth2 Config
	oauth2TwitchConfig = oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:8080/callback/twitch",
		Scopes:       []string{"user:read:email"},
		Endpoint:     twitch.Endpoint,
	}
}

type contextKey string

// NewAuthService initializes AuthService with dependencies
func NewAuthService(redisClient *redis.RedisClient, postgresClient *postgresql.PostgresClient, config *config.AppConfig) *AuthService {
	return &AuthService{
		PostgresClient: postgresClient,
		RedisClient:    redisClient,
		Config:         config,
	}
}

func (auth *AuthService) GenerateAPIKey() string {
	for {
		randomBytes := make([]byte, 32)
		_, err := rand.Read(randomBytes)
		if err != nil {
			logger.Error.Fatalf("Failed to generate API key: %v", err)
		}
		hash := sha256.Sum256(randomBytes)
		apiKey := base64.StdEncoding.EncodeToString(hash[:])

		var existingUser models.User
		if err := auth.PostgresClient.DB.Where("api_key = ?", apiKey).First(&existingUser).Error; err == gorm.ErrRecordNotFound {
			return apiKey
		}
	}
}

func (auth *AuthService) GenerateTokens(username, role string) (string, string, error) {
	accessClaims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 1 month expiration
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(auth.Config.Server.JwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken := fmt.Sprintf("refresh_%s_%d", username, time.Now().UnixNano())

	err = auth.RedisClient.Set(context.Background(), fmt.Sprintf("users:%s:refresh_token", username), refreshToken, time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}


func AdminHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ClaimsContextKey).(*Claims)
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("Welcome, Admin %s!", claims.Username)})
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value(ClaimsContextKey).(*Claims)
    if !ok {
        fmt.Println("Claims not found in context") // Debug log
        http.Error(w, "Claims not found", http.StatusUnauthorized)
        return
    }

    fmt.Println("User authenticated:", claims.Username) // Debug log

    utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
        "message": fmt.Sprintf("Welcome, User %s!", claims.Username),
    })
}


func (auth *AuthService) RegenAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(ClaimsContextKey).(*Claims)
	if !ok || claims.Username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Find the user in the database
	var user models.User
	if err := auth.PostgresClient.DB.Where("username = ?", claims.Username).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	newAPIKey := auth.GenerateAPIKey()
	user.APIKey = newAPIKey

	// Save the updated API key in the database
	if err := auth.PostgresClient.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update API key", http.StatusInternalServerError)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"message": "API key regenerated successfully",
		"api_key": newAPIKey,
	})
}
