package models

import "gorm.io/gorm"

type Menu struct {
	BasicModel
	StoreId         uint64 `json:"store_id" validate:"required"`
	Name            string `json:"name" validate:"required" gorm:"default:null"`
	QrCodeSrc       string `json:"qr_code_src" gorm:"default:null"`
	MenuURL         string `json:"menu_url" gorm:"default:null"`
	ColorOnThePrint string `json:"color_on_the_print" validate:"required,lte=100" gorm:"default:null"`
}

func (cg *Menu) List(pagination Pagination) (*Pagination, error) {
	var menus []*Menu

	db := new(gorm.DB)

	db.Scopes(paginate(menus, &pagination, db)).Find(&menus)
	pagination.Rows = menus

	return &pagination, nil
}
