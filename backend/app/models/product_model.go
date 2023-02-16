package models

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

type SapoProductResp struct {
	Content     string       `json:"content" validate:"required"`
	Summary     string       `json:"summary" gorm:"default:null"`
	CreatedOn   time.Time    `json:"created_on" gorm:"default:null"`
	Alias       string       `json:"alias"`
	ProductId   uint64       `json:"id"`
	Images      ImageArray   `json:"images"`
	Options     OptionArray  `json:"options"`
	ProductType string       `json:"product_type" gorm:"default:null"`
	PublishedOn time.Time    `json:"published_on" gorm:"default:null"`
	Tags        string       `json:"tags" gorm:"default:null"`
	ProductName string       `json:"name" gorm:"column:product_name" validate:"required"`
	ModifiedOn  time.Time    `json:"modified_on" gorm:"default:null"`
	Variants    VariantArray `json:"variants"`
	Vendor      string       `json:"vendor" gorm:"default:null"`
}

// User struct to describe Product object.
type Product struct {
	BasicModel
	SapoProductResp
	ProductId      uint64    `json:"product_id"`
	Price          float64   `json:"price" validate:"required" gorm:"default:null"`
	StoreId        uint64    `json:"store_id" validate:"required"`
	UserAppTokenId uint64    `json:"user_app_token_id" gorm:"default:null"`
	IsChargeTax    int       `json:"is_charge_tax" validate:"eq=0|eq=1"`
	ProductStatus  string    `json:"product_status"`
	Gateway        string    `json:"gateway" validate:"required"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateProductBody struct {
	Product
	CollectionId uint64 `json:"collection_id"`
	MenuId       uint64 `json:"menu_id" validate:"required"`
}

type ProductDBForm interface {
	GetProduct() *Product
	GetProductNameAlias() string
}

func (p *CreateProductBody) GetProduct() *Product {
	return &p.Product
}

func (p *CreateProductBody) GetProductNameAlias() string {
	lowerCaseName := (strings.ToLower((strings.Trim(p.ProductName, " "))))
	return strings.ReplaceAll(lowerCaseName, " ", "_")
}

type VariantArray []Variant

func (sla *VariantArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla VariantArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type ImageArray []Image

func (sla *ImageArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ImageArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type OptionArray []Option

func (sla *OptionArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla OptionArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type Variant struct {
	Id                  int       `json:"id"`
	ProductId           uint64    `json:"product_id"`
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
	Id        uint64   `json:"id"`
	ProductId uint64   `json:"product_id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Values    []string `json:"values"`
}

type Image struct {
	Id         uint64    `json:"id" gorm:"default:null"`
	ProductId  uint64    `json:"product_id"`
	Position   int       `json:"position" gorm:"default:null"`
	CreatedOn  time.Time `json:"created_on" gorm:"default:null"`
	ModifiedOn time.Time `json:"modified_on" gorm:"default:null"`
	Src        string    `json:"src"`
	VariantIds []int     `json:"variant_ids" gorm:"default:null"`
}
