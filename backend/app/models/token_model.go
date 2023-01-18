package models

// Renew struct to describe refresh token object.
type Renew struct {
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}
