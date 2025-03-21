package auth

import (
	"IcyAPI/internal/models"
	"IcyAPI/internal/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (auth *AuthService) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var signupRequest UserSignupRequest
	if err := json.NewDecoder(r.Body).Decode(&signupRequest); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, ErrorResponse{
			Message: "Bad request - Invalid input",
			Code:    "400_BAD_REQUEST",
		})
		return
	}

	if adminWhitelist[signupRequest.Username] {
		utils.WriteJSONResponse(w, http.StatusForbidden, ErrorResponse{
			Message: "Cannot create an admin account",
			Code:    "401_UNAUTHORIZED",
		})
		return
	}

	var user models.User
	if err := auth.PostgresClient.DB.Where("username = ?", signupRequest.Username).First(&user).Error; err == nil {
		utils.WriteJSONResponse(w, http.StatusConflict, ErrorResponse{
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
		utils.WriteJSONResponse(w, http.StatusInternalServerError, ErrorResponse{
			Message: "Failed to create user",
			Code:    "500_INTERNAL_SERVER_ERROR",
		})
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, SuccessResponse{
		Message: "User registered successfully",
	})
}

// LoginHandler handles user login
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
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

func (auth *AuthService) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.FormValue("refresh_token")
	username, err := auth.RedisClient.Client.Get(context.Background(), fmt.Sprintf("users:%s:refresh_token", refreshToken)).Result()

	if err != nil {
		http.Error(w, "Error checking refresh token", http.StatusInternalServerError)
		return
	}

	accessToken, _, _ := auth.GenerateTokens(username, "user")
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"access_token": accessToken})
}

// LogoutHandler invalidates the user's session by deleting the refresh token
func (auth *AuthService) LogoutHandler(w http.ResponseWriter, r *http.Request) {
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
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}
