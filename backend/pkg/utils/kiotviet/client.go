package kiotviet

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

const API_URL = "https://public.kiotapi.com"

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
		var connectToken ConnectTokenResponse

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

func ProductList(userId uint64, pageSize, currentItem int) (ProductsResponse, error) {
	var productsResponse ProductsResponse
	fmt.Println(userId)
	var app models.App

	errApp := database.Database.Joins("left join user_app_tokens on apps.id = user_app_tokens.app_id").Where("user_id = ?", userId).Where("gateway = ?", "kiotviet").First(&app)

	if errApp.Error != nil {
		return productsResponse, errApp.Error
	}

	userAppToken := new(models.UserAppToken)

	errUserAppToken := database.Database.Where("user_id = ?", userId).Where("app_id = ?", app.ID).First(&userAppToken)
	if errUserAppToken.Error != nil {
		return productsResponse, errUserAppToken.Error
	}

	a := fiber.AcquireAgent()
	a.Add("Accept", "application/json")
	a.Add("Retailer", app.AppName)
	a.Add("Authorization", "Bearer " + userAppToken.AccessToken)

	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(API_URL + fmt.Sprintf("/products?pageSize=%d&currentItem=%d", pageSize, currentItem))

	if err := a.Parse(); err != nil {
		return productsResponse, nil
	}

	code, body, errs := a.String()

	if (code == fiber.StatusUnauthorized) {
		connectToken, errConnectToken := ConnectToken(userId)
		if (connectToken && errConnectToken == nil) {
			return ProductList(userId, pageSize, currentItem)
		}
	}

	if (len(errs) != 0) {
		return productsResponse, errs[0]
	}

	if (code == fiber.StatusOK) {
		resErr := json.Unmarshal([]byte(body), &productsResponse)
		fmt.Println(resErr)
		if resErr != nil {
			return productsResponse, resErr
		}
	}

	return productsResponse, nil
}

func CollectionList(userId uint64, pageSize, currentItem int) (CollectionsResponse, error) {
	var collectionsResponse CollectionsResponse
	fmt.Println(userId)
	var app models.App

	errApp := database.Database.Joins("left join user_app_tokens on apps.id = user_app_tokens.app_id").Where("user_id = ?", userId).Where("gateway = ?", "kiotviet").First(&app)

	if errApp.Error != nil {
		return collectionsResponse, errApp.Error
	}

	userAppToken := new(models.UserAppToken)

	errUserAppToken := database.Database.Where("user_id = ?", userId).Where("app_id = ?", app.ID).First(&userAppToken)
	if errUserAppToken.Error != nil {
		return collectionsResponse, errUserAppToken.Error
	}

	a := fiber.AcquireAgent()
	a.Add("Accept", "application/json")
	a.Add("Retailer", app.AppName)
	a.Add("Authorization", "Bearer " + userAppToken.AccessToken)

	req := a.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(API_URL + fmt.Sprintf("/categories?pageSize=%d&currentItem=%d", pageSize, currentItem))

	if err := a.Parse(); err != nil {
		return collectionsResponse, nil
	}

	code, body, errs := a.String()

	if (code == fiber.StatusUnauthorized) {
		connectToken, errConnectToken := ConnectToken(userId)
		if (connectToken && errConnectToken == nil) {
			return CollectionList(userId, pageSize, currentItem)
		}
	}

	if (len(errs) != 0) {
		return collectionsResponse, errs[0]
	}

	if (code == fiber.StatusOK) {
		resErr := json.Unmarshal([]byte(body), &collectionsResponse)
		fmt.Println(resErr)
		if resErr != nil {
			return collectionsResponse, resErr
		}
	}

	return collectionsResponse, nil
}
