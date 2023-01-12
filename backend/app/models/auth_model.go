package models

import "github.com/golang-jwt/jwt/v4"

// SignUp struct to describe register a new user.
type SignUp struct {
	Email       string `json:"email" validate:"required,email,lte=255"`
	Password    string `json:"password" validate:"required,lte=255"`
	PhoneNumber string `json:"phone_number" validate:"lte=255"`
	UserRole    string `json:"user_role" validate:"required,lte=25"`
}

// SignIn struct to describe login user.
type SignIn struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,lte=255"`
}

// Google SignIn struct
type GoogleSignIn struct {
	GoogleJWT string `json:"code"`
}

type GoogleCallback struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Picture       string `json:"picture"`
	jwt.StandardClaims
}

type CreatePasswordClaims struct {
	Password string `json:"password" validate:"required,lte=255"`
}

type EmailResetPassword struct {
	Email string `json:"email" validate:"required,email,lte=255"`
}
