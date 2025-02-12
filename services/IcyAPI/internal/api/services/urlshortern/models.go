package urlshortern

type URLStore interface {
	Save(originalURL, shortURL string) (string, error)
	Get(shortURL string) (string, bool)
}

type URLMapping struct {
	ID          uint   `gorm:"primaryKey"`
	OriginalURL string `gorm:"unique;not null"`
	ShortURL    string `gorm:"unique;not null"`
}
