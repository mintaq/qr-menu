package sapo

import (
	"fmt"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func GetAuthURLByStore(store string) (string, error) {
	var app models.App

	tx := database.Database.First(&app, "gateway = ?", "sapo")
	if tx.Error != nil {
		return "", tx.Error
	}

	return fmt.Sprintf("https://%s/admin/oauth/authorize?client_id=%s&scope=%s&redirect_uri=%s", store, app.ApiKey, app.Scopes, app.RedirectUrl+"?store="+store), nil
}
