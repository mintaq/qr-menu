package models

// User struct to describe UserApp object.
type UserApp struct {
	ID          uint64 `json:"id" gorm:"primarykey"`
	UserId      uint64 `json:"user_id" validate:"required"`
	AppId       uint64 `json:"app_id" validate:"required"`
	AccessToken string `json:"access_token" validate:"required"`
	TimeModel
}

type Tabler interface {
	TableName() string
}

func (UserApp) TableName() string {
	return "user_app"
}
