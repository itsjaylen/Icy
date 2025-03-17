package auth

import (
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"
	"IcyAPI/internal/models"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	logger "itsjaylen/IcyLogger"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/twitch"
	"gorm.io/gorm"
)

var (
	JwtSecret          = []byte("SsYMx7hdNterwN011bzykWrMxjymmiu6")
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

// UserSignupRequest represents the user signup request body
type UserSignupRequest struct {
	Username string `json:"username" example:"newuser" binding:"required"`           // Username of the user
	Password string `json:"password" example:"securepassword123" binding:"required"` // Password for the user
}

// ErrorResponse represents a standard error response for bad requests
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// SuccessResponse represents a successful response for user signup
type SuccessResponse struct {
	Message string `json:"message"`
}

// LoginResponse represents the response schema for successful login.
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

// AuthService handles user-related endpoints
type AuthService struct {
	PostgresClient *postgresql.PostgresClient
	RedisClient    *redis.RedisClient
}

// NewAuthService initializes AuthService with dependencies
func NewAuthService(redisClient *redis.RedisClient, postgresClient *postgresql.PostgresClient) *AuthService {
	return &AuthService{
		PostgresClient: postgresClient,
		RedisClient:    redisClient,
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(JwtSecret)

	refreshToken := fmt.Sprintf("refresh_%s_%d", username, time.Now().UnixNano())
	err := auth.RedisClient.Set(context.Background(), fmt.Sprintf("users:%s:refresh_token", username), refreshToken, time.Hour)
	if err != nil {
		return "", "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (auth *AuthService) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var signupRequest UserSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&signupRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Bad request - Invalid input",
			Code:    "400_BAD_REQUEST",
		})
		return
	}

	if adminWhitelist[signupRequest.Username] {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Cannot create an admin account",
			Code:    "401_UNAUTHORIZED",
		})
		return
	}

	var user models.User
	if err := auth.PostgresClient.DB.Where("username = ?", signupRequest.Username).First(&user).Error; err == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Username already taken",
			Code:    "409_CONFLICT",
		})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(signupRequest.Password), bcrypt.DefaultCost)
	apiKey := auth.GenerateAPIKey()

	newUser := models.User{
		Username: signupRequest.Username,
		Password: string(hashedPassword),
		Role:     "user",
		APIKey:   apiKey,
	}
	if err := auth.PostgresClient.DB.Create(&newUser).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Message: "Failed to create user",
			Code:    "500_INTERNAL_SERVER_ERROR",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(SuccessResponse{
		Message: "User registered successfully",
	})
}

func (auth *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	username, password := r.FormValue("username"), r.FormValue("password")
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := auth.PostgresClient.DB.Where("username = ?", username).First(&user).Error; err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if user.Locked {
		if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
			http.Error(w, "Account locked. Try again later", http.StatusForbidden)
			return
		} else {
			user.Locked = false
			user.FailedAttempts = 0
			user.LockedUntil = nil
			auth.PostgresClient.DB.Save(&user)
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		user.FailedAttempts++
		if user.FailedAttempts >= 3 {
			user.Locked = true
			lockDuration := 15 * time.Minute
			user.LockedUntil = &time.Time{}
			*user.LockedUntil = time.Now().Add(lockDuration)
		}
		auth.PostgresClient.DB.Save(&user)

		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, _ := auth.GenerateTokens(username, user.Role)
	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		APIKey:       user.APIKey,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (auth *AuthService) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.FormValue("refresh_token")
	username, err := auth.RedisClient.Client.Get(context.Background(), fmt.Sprintf("users:%s:refresh_token", refreshToken)).Result()

	if err != nil {
		http.Error(w, "Error checking refresh token", http.StatusInternalServerError)
		return
	}

	accessToken, _, _ := auth.GenerateTokens(username, "user")
	json.NewEncoder(w).Encode(map[string]string{"access_token": accessToken})
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(ClaimsContextKey).(*Claims)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Welcome, Admin %s!", claims.Username)})
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(ClaimsContextKey).(*Claims)
	if !ok {
		http.Error(w, "Claims not found", http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Welcome, User %s!", claims.Username)})
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

	// Generate a new API key
	newAPIKey := auth.GenerateAPIKey()
	user.APIKey = newAPIKey

	// Save the updated API key in the database
	if err := auth.PostgresClient.DB.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update API key", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "API key regenerated successfully",
		"api_key": newAPIKey,
	})
}

// LogoutHandler invalidates the user's session by deleting the refresh token
func (auth *AuthService) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	// Bearer token format check
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	refreshToken := tokenParts[1]

	// Extract claims from token
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	// Delete refresh token from Redis (if stored there)
	err = auth.RedisClient.Delete(context.Background(), "refresh:"+claims.Username)
	if err != nil {
		http.Error(w, "Failed to revoke session", http.StatusInternalServerError)
		return
	}

	// Response confirming logout
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}