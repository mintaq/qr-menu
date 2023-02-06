package models

// User struct to describe App object.
type App struct {
	ID          uint64 `json:"id" gorm:"primarykey"`
	AppName     string `json:"app_name" gorm:"unique" validate:"required,lte=255"`
	ApiKey      string `json:"api_key" validate:"required"`
	SecretKey   string `json:"secret_key" validate:"required"`
	Scopes      string `json:"scopes"`
	RedirectUrl string `json:"redirect_url"`
	Gateway     string `json:"gateway"`
	TimeModel
}
