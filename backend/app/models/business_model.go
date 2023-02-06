package models

// User struct to describe Business object.
type Business struct {
	ID      uint64 `json:"id" gorm:"primarykey"`
	UserId  uint64 `json:"user_id"`
	AppId   uint64 `json:"app_id"`
	Country string `json:"country"`
	City    string `json:"city"`
	Address string `json:"address"`
	TimeModel
}
