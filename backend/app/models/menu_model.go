package models

import "gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"

type Menu struct {
	BasicModel
	StoreId uint64 `json:"store_id" validate:"required"`
	Name    string `json:"name" validate:"required" gorm:"default:null"`
	Role    string `json:"role" validate:"required,lte=25,oneof=main unpublished"`
}

func NewDefaultMenu(storeId uint64) *Menu {
	menu := new(Menu)
	menu.StoreId = storeId
	menu.Name = "Default"
	menu.Role = repository.MENU_ROLE_MAIN

	return menu
}
