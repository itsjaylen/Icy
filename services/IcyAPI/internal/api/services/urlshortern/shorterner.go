package urlshortern

import (
	"crypto/sha1"
)

const base62Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenerateShortURL generates a short URL from the original URL
func GenerateShortURL(originalURL string) string {
	hash := sha1.New()
	hash.Write([]byte(originalURL))
	hashBytes := hash.Sum(nil)

	shortURL := ""
	for i := 0; i < len(hashBytes); i++ {
		shortURL += string(base62Charset[int(hashBytes[i])%len(base62Charset)])
	}
	return shortURL[:8] 
}