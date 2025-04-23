package urlshortern

import (
	"fmt"
	"net/http"
)

func HandleShortenURL(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originalURL := r.URL.Query().Get("url")
		if originalURL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		shortURL := GenerateShortURL(originalURL)

		// Try to store mapping, handle duplicate short URL
		createdShortURL, err := store.Save(originalURL, shortURL)
		if err != nil {
			http.Error(w, "Error saving URL", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Shortened URL: http://localhost:9800/%s", createdShortURL)
	}
}

func HandleRedirect(store URLStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortURL := r.URL.Path[1:]

		originalURL, exists := store.Get(shortURL)
		if !exists {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		// Send the JavaScript-based redirect
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<title>Redirecting...</title>
			</head>
			<body>
				<p>Redirecting to <a href="%s">%s</a></p>
				<script type="text/javascript">
					// Redirect after a short delay to keep the short URL in the address bar
					window.location.replace("%s");
				</script>
			</body>
			</html>
		`, originalURL, originalURL, originalURL)
	}
}

