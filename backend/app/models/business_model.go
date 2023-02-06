package models

// User struct to describe Business object.
type Business struct {
	ID        uint64 `json:"id" gorm:"primarykey"`
	UserId    uint64 `json:"user_id" validate:"required"`
	StoreName string `json:"store_name" validate:"required,lte=255"`
	StoreURL  string `json:"store_url" validate:"required"`
	Country   string `json:"country" validate:"required,lte=255"`
	City      string `json:"city" validate:"required,lte=255"`
	Address   string `json:"address" validate:"required"`
	TimeModel
}
