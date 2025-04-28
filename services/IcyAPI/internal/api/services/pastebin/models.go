package pastebin

import "time"

type Paste struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string    `gorm:"type:text"`
	Content     string    `gorm:"type:text;not null"`
	Syntax      string    `gorm:"type:text;default:'plaintext'"`
	ImageURL    string    `gorm:"type:text"`
	Views       int       `gorm:"default:0"`
	DeleteToken string    `gorm:"type:uuid;not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
