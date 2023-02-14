package kiotviet

import (
	"fmt"
	"strconv"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func GetAuthURLByStore(store string, userId uint64) (string, error) {
	var app models.App

	tx := database.Database.First(&app, "gateway = ?", repository.GATEWAY_SAPO)
	if tx.Error != nil {
		return "", tx.Error
	}

	redirectURL := app.RedirectUrl + "?user_id=" + strconv.FormatUint(userId, 10)

	return fmt.Sprintf("https://%s/admin/oauth/authorize?client_id=%s&scope=%s&redirect_uri=%s", store, app.ApiKey, app.Scopes, redirectURL), nil
}
