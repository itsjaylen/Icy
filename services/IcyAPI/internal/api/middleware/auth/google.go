package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/itsjaylen/IcyAPI/internal/models"
	"golang.org/x/oauth2"
)

// GoogleLoginHandler redirects the user to Google's OAuth consent page.
func GoogleLoginHandler(writer http.ResponseWriter, request *http.Request) {
	url := oauth2GoogleConfig.AuthCodeURL(oauth2StateString, oauth2.AccessTypeOffline)
	http.Redirect(writer, request, url, http.StatusFound)
}

// GoogleCallbackHandler handles the OAuth2 callback and retrieves the user's Google profile.
func (auth *Service) GoogleCallbackHandler(writer http.ResponseWriter, request *http.Request) {
	// Exchange the authorization code for an access token
	code := request.URL.Query().Get("code")
	if code == "" {
		http.Error(writer, "Code not found", http.StatusBadRequest)

		return
	}

	token, err := oauth2GoogleConfig.Exchange(request.Context(), code)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)

		return
	}

	// Get user information from Google
	client := oauth2GoogleConfig.Client(request.Context(), token)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
defer cancel()

req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", nil)
if err != nil {
	http.Error(writer, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
	return
}

resp, err := client.Do(req)
if err != nil {
	http.Error(writer, fmt.Sprintf("Failed to fetch user info: %v", err), http.StatusInternalServerError)
	return
}
defer resp.Body.Close()

	var userInfo struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		FullName string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(writer, "Failed to decode user info", http.StatusInternalServerError)

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
			http.Error(writer, "Failed to create user", http.StatusInternalServerError)

			return
		}
	}

	// Generate JWT tokens
	accessToken, refreshToken, _ := auth.GenerateTokens(user.Username, user.Role)

	// Send the tokens back to the user
	json.NewEncoder(writer).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"api_key":       user.APIKey,
	})
}
