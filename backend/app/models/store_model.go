package models

import "regexp"

// User struct to describe Store object.
type Store struct {
	BasicModel
	UserId    uint64 `json:"user_id" validate:"required"`
	Name      string `json:"name" validate:"required,lte=255"`
	Subdomain string `json:"subdomain" validate:"required"`
	Country   string `json:"country" validate:"lte=255"`
	City      string `json:"city" validate:"lte=255"`
	Address   string `json:"address" validate:""`
}

func (s *Store) VerifySubdomain() (bool, error) {
	return regexp.Match(`^[a-z0-9_-]+\.dingdoong\.io$`, []byte(s.Subdomain))
}
