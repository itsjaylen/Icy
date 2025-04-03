package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/itsjaylen/IcyAPI/internal/models"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

// GitHub OAuth Handlers.
func GithubLoginHandler(writer http.ResponseWriter, request *http.Request) {
	url := oauth2GithubConfig.AuthCodeURL(oauth2StateString, oauth2.AccessTypeOffline)
	http.Redirect(writer, request, url, http.StatusFound)
}

func (auth *Service) GithubCallbackHandler(writer http.ResponseWriter, request *http.Request) {
	code := request.URL.Query().Get("code")
	if code == "" {
		http.Error(writer, "Code not found", http.StatusBadRequest)

		return
	}

	token, err := oauth2GithubConfig.Exchange(request.Context(), code)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to exchange token: %v", err), http.StatusInternalServerError)

		return
	}

	client := oauth2GithubConfig.Client(request.Context(), token)

	// Fetch user info from GitHub
	userResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(writer, fmt.Sprintf("Failed to get user info: %v", err), http.StatusInternalServerError)

		return
	}
	defer userResp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
		ID    int    `json:"id"`
	}

	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		http.Error(writer, "Failed to decode user info", http.StatusInternalServerError)

		return
	}

	// If GitHub doesn't return an email, request emails explicitly
	if userInfo.Email == "" {
		emailResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			http.Error(writer, "Failed to fetch emails", http.StatusInternalServerError)

			return
		}
		defer emailResp.Body.Close()

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
			http.Error(writer, "Failed to decode email response", http.StatusInternalServerError)

			return
		}

		// Pick the first verified primary email
		for _, e := range emails {
			if e.Primary && e.Verified {
				userInfo.Email = e.Email

				break
			}
		}

		// If still no email, return an error
		if userInfo.Email == "" {
			http.Error(writer, "No verified email found for this GitHub account", http.StatusBadRequest)

			return
		}
	}

	// Check if the user exists in the database
	var user models.User
	err = auth.PostgresClient.DB.Where("username = ?", userInfo.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// User not found, create a new one
			user = models.User{
				Username: userInfo.Email,
				Role:     "user",
				APIKey:   auth.GenerateAPIKey(),
			}
			if err := auth.PostgresClient.DB.Create(&user).Error; err != nil {
				http.Error(writer, "Failed to create user", http.StatusInternalServerError)

				return
			}
		} else {
			http.Error(writer, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)

			return
		}
	}

	// Generate JWT tokens
	accessToken, refreshToken, _ := auth.GenerateTokens(user.Username, user.Role)

	// Send the tokens back to the user
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"api_key":       user.APIKey,
	})
}
