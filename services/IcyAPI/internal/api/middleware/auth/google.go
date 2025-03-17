package auth

import (
	"IcyAPI/internal/models"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

// GoogleLoginHandler redirects the user to Google's OAuth consent page
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauth2GoogleConfig.AuthCodeURL(oauth2StateString, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

// GoogleCallbackHandler handles the OAuth2 callback and retrieves the user's Google profile
func (auth *AuthService) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	// Exchange the authorization code for an access token
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := oauth2GoogleConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	// Get user information from Google
	client := oauth2GoogleConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		FullName string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	// Check if the user exists in the database
	var user models.User
	if err := auth.PostgresClient.DB.Where("username = ?", userInfo.Email).First(&user).Error; err != nil {
		// User not found, create new user
		user = models.User{
			Username: userInfo.Email,
			Role:     "user",
			APIKey:   auth.GenerateAPIKey(),
		}
		if err := auth.PostgresClient.DB.Create(&user).Error; err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
	}

	// Generate JWT tokens
	accessToken, refreshToken, _ := auth.GenerateTokens(user.Username, user.Role)

	// Send the tokens back to the user
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"api_key":       user.APIKey,
	})
}