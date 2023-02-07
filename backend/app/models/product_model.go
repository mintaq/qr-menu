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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
