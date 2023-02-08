package models

// User struct to describe User object.
type User struct {
	BasicModel
	Email        string `json:"email" gorm:"unique" validate:"required,email,lte=255"`
	PasswordHash string `json:"password_hash,omitempty" validate:"lte=255"`
	PhoneNumber  string `json:"phone_number" validate:"lte=255"`
	UserStatus   int    `json:"user_status" validate:"required,len=1"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	UserRole     string `json:"user_role" gorm:"default:user" validate:"required,lte=25"`
	UserImage    string `json:"user_image"`
}
