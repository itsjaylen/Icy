package auth

import (
	"IcyAPI/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// TwitchLoginHandler redirects the user to Twitch's OAuth consent page
func TwitchLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauth2TwitchConfig.AuthCodeURL(oauth2StateString, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (auth *AuthService) TwitchCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for an access token
	token, err := oauth2TwitchConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch user information from Twitch
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users", nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Client-ID", oauth2TwitchConfig.ClientID)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var twitchResponse struct {
		Data []struct {
			ID          string `json:"id"`
			Login       string `json:"login"`
			DisplayName string `json:"display_name"`
			Email       string `json:"email"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&twitchResponse); err != nil {
		http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
		return
	}

	if len(twitchResponse.Data) == 0 {
		http.Error(w, "No user data found", http.StatusInternalServerError)
		return
	}

	twitchUser := twitchResponse.Data[0]

	// Check if the user exists in the database
	var user models.User
	err = auth.PostgresClient.DB.Where("username = ?", twitchUser.Login).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found, create a new one
			user = models.User{
				Username: twitchUser.Login,
				Role:     "user",
				APIKey:   auth.GenerateAPIKey(),
			}
			if err := auth.PostgresClient.DB.Create(&user).Error; err != nil {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Generate JWT tokens
	accessToken, refreshToken, _ := auth.GenerateTokens(user.Username, user.Role)

	// Send the tokens back to the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"api_key":       user.APIKey,
	})
}