package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

type Collection struct {
	BasicModel
	SapoCollectionResp
	CollectionId   uint64 `json:"collection_id"`
	StoreId        uint64 `json:"store_id" validate:"required"`
	UserAppTokenId uint64 `json:"user_app_token_id" gorm:"default:null"`
	Gateway        string `json:"gateway"`
}

type SapoCollectionResp struct {
	CollectionId uint64          `json:"id"`
	Description  string          `json:"description" gorm:"default:null"`
	Alias        string          `json:"alias" gorm:"default:null"`
	Name         string          `json:"name" validate:"required"`
	Image        CollectionImage `json:"image" gorm:"default:null"`
}

type CollectionImage struct {
	Id        uint64    `json:"id"`
	Src       string    `json:"src"`
	CreatedOn time.Time `json:"created_on" gorm:"default:null"`
}

func (sla *CollectionImage) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CollectionImage) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

func (p *Collection) GetNameAlias() string {
	lowerCaseName := (strings.ToLower((strings.Trim(p.Name, " "))))
	return strings.ReplaceAll(lowerCaseName, " ", "_")
}
