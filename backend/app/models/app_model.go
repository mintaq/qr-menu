package models

// User struct to describe App object.
type App struct {
	ID          int    `json:"id" gorm:"primarykey"`
	AppName     string `json:"app_name" gorm:"unique" validate:"required,lte=255"`
	ApiKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	Scopes      int    `json:"scopes"`
	RedirectUrl string `json:"redirect_url"`
	Gateway     string `json:"gateway"`
	TimeModel
}
