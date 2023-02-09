package models

type Collect struct {
	BasicModel
	CollectionId uint64 `json:"collection_id"`
	ProductId    uint64 `json:"product_id"`
	StoreId      uint64 `json:"store_id"`
	Position     int    `json:"position"`
}
