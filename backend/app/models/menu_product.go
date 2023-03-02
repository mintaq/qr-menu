package models

type MenuProduct struct {
	BasicModel
	StoreId   uint64 `json:"store_id" validate:"required"`
	MenuId    uint64 `json:"menu_id" validate:"required"`
	ProductId uint64 `json:"product_id" validate:"required"`
}

func (MenuProduct) TableName() string {
	return "menu_product"
}
