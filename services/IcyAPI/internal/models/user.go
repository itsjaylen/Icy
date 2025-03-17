package models

import "time"

// User model for PostgreSQL
type User struct {
	ID             uint       `json:"id" gorm:"primaryKey" example:"1"`
	Username       string     `json:"username" gorm:"unique;not null" example:"newuser"`
	Password       string     `json:"password" example:"hashedpassword123"`
	Role           string     `json:"role" example:"user"`
	APIKey         string     `json:"api_key" gorm:"unique" example:"abcd1234"`
	Locked         bool       `json:"locked" gorm:"default:false"`
	FailedAttempts int        `json:"failed_attempts" gorm:"default:0"`
	LockedUntil    *time.Time `json:"locked_until,omitempty"`
	CreatedAt      time.Time  `json:"created_at" example:"2025-02-28T15:00:00Z"`
	UpdatedAt      time.Time  `json:"updated_at" example:"2025-02-28T15:00:00Z"`
}

// UserSignupRequest represents the user signup request body
type UserSignupRequest struct {
	Username string `json:"username" example:"newuser" binding:"required"`           // Username of the user
	Password string `json:"password" example:"securepassword123" binding:"required"` // Password for the user
}
