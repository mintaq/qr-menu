package models

import (
	"time"
)

type BasicModel struct {
	ID        uint64    `gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // Set to current time if it is zero on creating
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // Set to current unix seconds on updating or if it is zero on creating
}
