package models

type Site struct {
	BasicModel
	UserId    uint64 `json:"user_id"`
	Name      string `json:"name" validate:"required"`
	Subdomain string `json:"subdomain" validate:"required"`
}
