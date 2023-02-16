package models

type MenuProduct struct {
	BasicModel
	StoreId   uint64 `json:"store_id"`
	MenuId    uint64 `json:"menu_id"`
	ProductId uint64 `json:"product_id"`
}
