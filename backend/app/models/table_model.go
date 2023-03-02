package models

import (
	"fmt"
	"os"
)

type Table struct {
	BasicModel
	StoreId         uint64 `json:"store_id" validate:"required"`
	Name            string `json:"name" validate:"required" gorm:"default:null"`
	QrCodeSrc       string `json:"qr_code_src" gorm:"default:null"`
	TableURL        string `json:"table_url" gorm:"default:null"`
	ColorOnThePrint string `json:"color_on_the_print" validate:"required,lte=100" gorm:"default:null"`
}

func (t *Table) GetTableQrImageSrc() string {
	return fmt.Sprintf("%s%d.png", os.Getenv("QR_CODE_TABLE_IMAGE_PREFIX"), t.ID)
}

func (t *Table) GetTableURL(subdomain string, menuId int) string {
	return fmt.Sprintf("https://%s?menu=%d&table=%d", subdomain, menuId, t.ID)
}
