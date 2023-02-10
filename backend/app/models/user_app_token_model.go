package models

// User struct to describe UserAppToken object.
type UserAppToken struct {
	BasicModel
	UserId      uint64 `json:"user_id" validate:"required"`
	AppId       uint64 `json:"app_id" validate:"required"`
	StoreDomain string `json:"store_domain" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
}
