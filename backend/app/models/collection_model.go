package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Collection struct {
	BasicModel
	SapoCollectionResp
	CollectionId   uint64    `json:"collection_id"`
	StoreId        uint64    `json:"store_id" validate:"required"`
	UserAppTokenId uint64    `json:"user_app_token_id" gorm:"default:null"`
	Gateway        string    `json:"gateway"`
	Products       []Product `json:"products" gorm:"-"`
}

type CollectionWithProducts struct {
	Collection `json:"collection"`
	Products   []Product `json:"products" gorm:"-"`
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

type CollectionQueryParams struct {
	PaginationQueryParams
	Name     string   `query:"name"`
	StoreId  uint64   `query:"store_id" validate:"required"`
	Includes []string `query:"includes"`
}

func (c *Collection) GetProducts(db *gorm.DB) ([]Product, error) {
	products := []Product{}

	if tx := db.Model(Product{}).Joins("left join collects on collects.product_id = products.product_id and collects.collection_id = ?", c.CollectionId).Where("products.store_id = ? AND collects.collection_id = ?", c.StoreId, c.CollectionId).Find(&products); tx.Error != nil {
		return nil, tx.Error
	}

	return products, nil
}
