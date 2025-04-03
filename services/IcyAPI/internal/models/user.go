package models

import "time"

// User model for PostgreSQL.
type User struct {
	CreatedAt      time.Time  `json:"createdAt" example:"2025-02-28T15:00:00Z"`
	UpdatedAt      time.Time  `json:"updatedAt" example:"2025-02-28T15:00:00Z"`
	LockedUntil    *time.Time `json:"lockedUntil,omitempty"`
	Username       string     `json:"username" gorm:"unique;not null" example:"newuser"`
	Password       string     `json:"password" example:"hashedpassword123"`
	Role           string     `json:"role" example:"user"`
	APIKey         string     `json:"apiKey" gorm:"unique" example:"abcd1234"`
	ID             uint       `json:"id" gorm:"primaryKey" example:"1"`
	FailedAttempts int        `json:"failedAttempts" gorm:"default:0"`
	Locked         bool       `json:"locked" gorm:"default:false"`
}

// UserSignupRequest represents the user signup request body.
type UserSignupRequest struct {
	Username string `json:"username" example:"newuser" binding:"required"`
	Password string `json:"password" example:"securepassword123" binding:"required"`
}
