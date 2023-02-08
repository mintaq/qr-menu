package models

import "time"

type Product struct {
	BasicModel
	UserAppId   uint64    `json:"user_app_id" validate:"required"`
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	CreatedOn   time.Time `json:"created_on"`
	Alias       string    `json:"alias"`
	ProductId   uint64    `json:"product_id" validate:"required"`
	Images      []Image   `json:"images"`
	Options     []Option  `json:"options"`
	ProductType string    `json:"product_type" validate:"required"`
	PublishedOn time.Time `json:"published_on"`
	Tags        string    `json:"tags"`
	ProductName string    `json:"product_name" validate:"required"`
	ModifiedOn  time.Time `json:"modified_on"`
	Variants    []Variant `json:"variants"`
	Vendor      string    `json:"vendor"`
	Gateway     string    `json:"gateway" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductResp struct {
	Content     string    `json:"content"`
	Summary     string    `json:"summary"`
	CreatedOn   time.Time `json:"created_on"`
	Alias       string    `json:"alias"`
	ProductId   uint64    `json:"product_id" validate:"required"`
	Images      string    `json:"images"`
	Options     string    `json:"options"`
	ProductType string    `json:"product_type" validate:"required"`
	PublishedOn time.Time `json:"published_on"`
	Tags        string    `json:"tags"`
	ProductName string    `json:"product_name" validate:"required"`
	ModifiedOn  time.Time `json:"modified_on"`
	Variants    string    `json:"variants"`
	Vendor      string    `json:"vendor"`
	Gateway     string    `json:"gateway" validate:"required"`
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
