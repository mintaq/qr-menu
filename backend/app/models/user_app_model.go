package models

// User struct to describe UserApp object.
type UserApp struct {
	ID     uint64 `json:"id" gorm:"primarykey"`
	UserId uint64 `json:"user_id"`
	AppId  uint64 `json:"app_id"`
	TimeModel
}
