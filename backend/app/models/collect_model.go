package models

type Collect struct {
	BasicModel
	CollectionId   uint64 `json:"collection_id"`
	ProductId      uint64 `json:"product_id"`
	UserAppTokenId uint64 `json:"user_app_token_id" validate:"required"`
	Position       int    `json:"position"`
}
