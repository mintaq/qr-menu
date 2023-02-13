package models

// User struct to describe Store object.
type Store struct {
	BasicModel
	UserId    uint64 `json:"user_id" validate:"required"`
	Name      string `json:"name" validate:"required,lte=255"`
	Subdomain string `json:"subdomain" validate:"required"`
	Country   string `json:"country" validate:"required,lte=255"`
	City      string `json:"city" validate:"required,lte=255"`
	Address   string `json:"address" validate:"required"`
}
