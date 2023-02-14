package models

import (
	"database/sql/driver"
	"encoding/json"
)

type Theme struct {
	BasicModel
	StoreId    uint64      `json:"store_id" validate:"required"`
	CoverImage string      `json:"cover_image" gorm:"default:null"`
	Colors     ThemeColors `json:"colors" validate:"required,json"`
}

type ThemeColors struct {
	Button       string `json:"button" validate:"required"`
	CategoryText string `json:"category_text" validate:"required"`
	ProductText  string `json:"product_text" validate:"required"`
	Background   string `json:"background" validate:"required"`
}

func (sla *ThemeColors) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ThemeColors) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}
