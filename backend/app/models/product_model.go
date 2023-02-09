package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type SapoProductResp struct {
	Content     string       `json:"content"`
	Summary     string       `json:"summary"`
	CreatedOn   time.Time    `json:"created_on"`
	Alias       string       `json:"alias"`
	ProductId   uint64       `json:"id"`
	Images      ImageArray   `json:"images"`
	Options     OptionArray  `json:"options"`
	ProductType string       `json:"product_type" validate:"required"`
	PublishedOn time.Time    `json:"published_on"`
	Tags        string       `json:"tags"`
	ProductName string       `json:"name" gorm:"column:product_name" validate:"required"`
	ModifiedOn  time.Time    `json:"modified_on"`
	Variants    VariantArray `json:"variants"`
	Vendor      string       `json:"vendor"`
}

// User struct to describe Product object.
type Product struct {
	BasicModel
	SapoProductResp
	ProductId uint64    `json:"product_id"`
	StoreId   uint64    `json:"store_id" validate:"required"`
	Gateway   string    `json:"gateway" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VariantArray []Variant

func (sla *VariantArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla VariantArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type ImageArray []Variant

func (sla *ImageArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ImageArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type OptionArray []Variant

func (sla *OptionArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla OptionArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type Variant struct {
	Id                  int       `json:"id"`
	ProductId           int       `json:"product_id"`
	Title               string    `json:"title"`
	Price               float32   `json:"price"`
	Sku                 string    `json:"sku"`
	Position            int       `json:"position"`
	Grams               int       `json:"grams"`
	InventoryManagement string    `json:"inventory_management"`
	Option1             string    `json:"option1"`
	Option2             string    `json:"option2"`
	Option3             string    `json:"option3"`
	CreatedOn           time.Time `json:"created_on"`
	ModifiedOn          time.Time `json:"modified_on"`
	RequiresShipping    bool      `json:"requires_shipping"`
	Barcode             string    `json:"barcode"`
	InventoryQuantity   int       `json:"inventory_quantity"`
	ImageId             int       `json:"image_id"`
	Weight              float32   `json:"weight"`
	WeightUnit          string    `json:"weight_unit"`
}

type Option struct {
	Id        int      `json:"id"`
	ProductId int      `json:"product_id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Values    []string `json:"values"`
}

type Image struct {
	Id         int       `json:"id"`
	ProductId  int       `json:"product_id"`
	Position   int       `json:"position"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
	Src        string    `json:"src"`
	VariantIds []int     `json:"variant_ids"`
}
