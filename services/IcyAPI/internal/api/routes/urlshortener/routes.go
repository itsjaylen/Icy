// routes.go (package urlshortener)

package urlshortener

import (
	"net/http"
	"time"

	"github.com/itsjaylen/IcyAPI/internal/api/middleware"
	"github.com/itsjaylen/IcyAPI/internal/api/services/urlshortern"
	"gorm.io/gorm"
)

// RegisterRoutes registers URL shortener-related routes.
func RegisterRoutes(mux *http.ServeMux, db *gorm.DB) {
	store := urlshortern.NewURLStore(db)

	mux.Handle("/shorten", middleware.RateLimitMiddleware(urlshortern.HandleShortenURL(store), 5*time.Second, 3))
	mux.Handle("/", middleware.RateLimitMiddleware(urlshortern.HandleRedirect(store), 5*time.Second, 3))
}
