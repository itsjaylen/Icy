package urlshortern

import (
	"gorm.io/gorm"
)


type URLStore interface {
	Save(originalURL, shortURL string) (string, error)
	Get(shortURL string) (string, bool)
}

type URLMapping struct {
	ID          uint   `gorm:"primaryKey"`
	OriginalURL string `gorm:"unique;not null"`
	ShortURL    string `gorm:"unique;not null"`
}

type PostgreSQLStore struct {
	db *gorm.DB
}

func NewURLStore(db *gorm.DB) *PostgreSQLStore {
	return &PostgreSQLStore{db: db}
}

// Save stores the URL mapping. If the short URL exists, it returns the existing short URL.
func (s *PostgreSQLStore) Save(originalURL, shortURL string) (string, error) {
	var existingURL URLMapping
	err := s.db.Where("short_url = ?", shortURL).First(&existingURL).Error
	if err == nil {
		return existingURL.ShortURL, nil
	}

	err = s.db.Create(&URLMapping{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}).Error
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *PostgreSQLStore) Get(shortURL string) (string, bool) {
	var urlMapping URLMapping
	err := s.db.Where("short_url = ?", shortURL).First(&urlMapping).Error
	if err != nil {
		return "", false
	}
	return urlMapping.OriginalURL, true
}



// ShortenAndSaveURL is a helper function that combines generating and saving the short URL
func ShortenAndSaveURL(store URLStore, originalURL string) (string, error) {
	// Generate the short URL
	shortURL := GenerateShortURL(originalURL)

	createdShortURL, err := store.Save(originalURL, shortURL)
	if err != nil {
		return "", err
	}

	return createdShortURL, nil
}
