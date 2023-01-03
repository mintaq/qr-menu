package models

// User struct to describe User object.
type User struct {
	ID           int    `json:"id" gorm:"primarykey"`
	Email        string `json:"email" gorm:"unique" validate:"required,email,lte=255"`
	PasswordHash string `json:"password_hash,omitempty" validate:"required,lte=255"`
	UserStatus   int    `json:"user_status" validate:"required,len=1"`
	UserRole     string `json:"user_role" validate:"required,lte=25"`
	TimeModel
}
