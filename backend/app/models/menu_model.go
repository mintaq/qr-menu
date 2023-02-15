package models

type Menu struct {
	BasicModel
	StoreId         uint64 `json:"store_id" validate:"required"`
	Name            string `json:"name" validate:"required" gorm:"default:null"`
	QrCodeSrc       string `json:"qr_code_src" gorm:"default:null"`
	MenuURL         string `json:"menu_url" gorm:"default:null"`
	ColorOnThePrint string `json:"color_on_the_print" validate:"required,lte=100" gorm:"default:null"`
}
