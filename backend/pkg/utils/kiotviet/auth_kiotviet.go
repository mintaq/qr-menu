package kiotviet

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

type RespConnectToken struct {
	AccessToken     string `json:"access_token"`
	ExpiresIn     int `json:"expires_in"`
	TokenType     string `json:"token_type"`
	Scope     string `json:"scope"`
}

func ConnectToken(userId uint64) (bool, error) {
	var app models.App

	appErr := database.Database.Joins("left join user_app_tokens on apps.id = user_app_tokens.app_id").Where("user_id = ?", userId).Where("gateway = ?", "kiotviet").First(&app)

	if appErr.Error != nil {
		return false, appErr.Error
	}

	args := fiber.AcquireArgs()
	args.Set("scopes", "PublicApi.Access")
	args.Set("grant_type", "client_credentials")
	args.Set("client_id", app.ApiKey)
	args.Set("client_secret", app.SecretKey)
	a := fiber.AcquireAgent()
	a.Add("Content-Type", "application/x-www-form-urlencoded")
	a.Form(args)

	req := a.Request()
	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI("https://id.kiotviet.vn/connect/token")

	if err := a.Parse(); err != nil {
		return false, nil
	}

	code, body, errs := a.String()

	if (len(errs) != 0) {
		return false, errs[0]
	}

	if (code == fiber.StatusOK) {
		var connectToken RespConnectToken

		resErr := json.Unmarshal([]byte(body), &connectToken)
		if resErr != nil {
			return false, resErr
		}

		userAppToken := new(models.UserAppToken)
		userAppTokenErr := database.Database.Where("user_id = ?", userId).Where("app_id = ?", app.ID).First(&userAppToken)

		if userAppTokenErr.Error != nil {
			return false, userAppTokenErr.Error
		}

		userAppToken.AccessToken = connectToken.AccessToken

		saveErr := database.Database.Save(userAppToken)
		if saveErr.Error != nil {
			return false, saveErr.Error
		}
	}

	return true, nil
}
