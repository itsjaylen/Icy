// routes.go (package urlshortener)

package urlshortener

import (
	"net/http"

	"github.com/itsjaylen/IcyAPI/internal/api/services/urlshortern"
	"gorm.io/gorm"
)

// RegisterRoutes registers URL shortener-related routes.
func RegisterRoutes(mux *http.ServeMux, db *gorm.DB) {
	store := urlshortern.NewURLStore(db)

	mux.Handle("/shorten", urlshortern.HandleShortenURL(store))
	mux.Handle("/", urlshortern.HandleRedirect(store))
}
