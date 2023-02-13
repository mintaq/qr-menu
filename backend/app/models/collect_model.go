package models

type Collect struct {
	BasicModel
	CollectionId   uint64 `json:"collection_id"`
	ProductId      uint64 `json:"product_id"`
	StoreId        uint64 `json:"store_id" validate:"required"`
	UserAppTokenId uint64 `json:"user_app_token_id" gorm:"default:null"`
	Position       int    `json:"position" gorm:"default:null"`
}
