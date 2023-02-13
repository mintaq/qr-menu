package models

type QrCode struct {
	BasicModel
	StoreId uint64 `json:"store_id"`
	Name    string `json:"name"`
	Image   string `json:"image"`
}
