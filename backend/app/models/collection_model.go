package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type Collection struct {
	BasicModel
	SapoCollectionResp
	CollectionId uint64 `json:"collection_id"`
	StoreId      uint64 `json:"store_id"`
	Gateway      string `json:"gateway"`
}

type SapoCollectionResp struct {
	CollectionId uint64      `json:"id"`
	Description  string      `json:"description"`
	Alias        string      `json:"alias"`
	Name         string      `json:"name"`
	Image        ImageObject `json:"image"`
}

type ImageObject struct {
	Src       string    `json:"src"`
	CreatedOn time.Time `json:"created_on"`
}

func (sla *ImageObject) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ImageObject) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}
