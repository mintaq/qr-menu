package models

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

// User struct to describe Store object.
type Store struct {
	BasicModel
	UserId    uint64 `json:"user_id" validate:"required"`
	Subdomain string `json:"subdomain" validate:"required"`
	StoreUpdatableData
}

type StoreUpdatableData struct {
	Name    string `json:"name" validate:"required,lte=255"`
	Country string `json:"country" validate:"lte=255"`
	City    string `json:"city" validate:"lte=255"`
	Address string `json:"address" validate:""`
}

func (s *Store) VerifySubdomain() (bool, error) {
	return regexp.Match(`^[a-z0-9_-]+\.dingdoong\.io$`, []byte(s.Subdomain))
}

func (s *Store) UpdateStore(ud *StoreUpdatableData) {
	if ud.Name != "" {
		s.Name = ud.Name
	}
	if ud.Country != "" {
		s.Country = ud.Country
	}
	if ud.City != "" {
		log.Println(`alo`)
		s.City = ud.City
	}
	if ud.Address != "" {
		s.Address = ud.Address
	}
}

func (s *Store) AddSuffixToSubdomain() {
	s.Subdomain = fmt.Sprintf("%s.%s", s.Subdomain, os.Getenv("SUBDOMAIN_SUFFIX"))
}
