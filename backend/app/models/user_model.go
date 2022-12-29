package models

import (
	"time"

	"gorm.io/gorm"
)

// User struct to describe User object.
type User struct {
	gorm.Model
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email" validate:"required,email,lte=255"`
	PasswordHash string    `json:"password_hash,omitempty" validate:"required,lte=255"`
	UserStatus   int       `json:"user_status" validate:"required,len=1"`
	UserRole     string    `json:"user_role" validate:"required,lte=25"`
}
