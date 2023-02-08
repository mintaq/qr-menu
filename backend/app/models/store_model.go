package models

// User struct to describe Store object.
type Store struct {
	BasicModel
	UserId      uint64 `json:"user_id" validate:"required"`
	AppId       uint64 `json:"app_id" validate:"required"`
	Store       string `json:"store" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
}
