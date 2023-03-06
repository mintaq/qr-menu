package kiotviet

import (
	"fmt"
	"log"
	"time"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"gorm.io/gorm/clause"
)

func SyncProducts(userId uint64, storeId uint64, pageSize int, currentItem int) (int, error) {
	log.Println("SyncProducts: Processing...")

	store := new(models.Store)

	errStore := database.Database.Where("user_id = ?", userId).Where("id = ?", storeId).First(&store)
	if errStore.Error != nil {
		return 0, errStore.Error
	}

	var app models.App
	errApp := database.Database.Joins("left join user_app_tokens on apps.id = user_app_tokens.app_id").Where("user_id = ?", userId).Where("gateway = ?", "kiotviet").First(&app)

	if errApp.Error != nil {
		return 0, errApp.Error
	}

	var userAppToken models.UserAppToken

	errUserAppToken := database.Database.Where("user_id = ?", userId).Where("app_id = ?", app.ID).First(&userAppToken)
	if errUserAppToken.Error != nil {
		return 0, errUserAppToken.Error
	}

	var productsResponse ProductsResponse

	productsResponse, errProductsResponse := ProductList(userId, pageSize, currentItem)
	if (errProductsResponse != nil) {
		return 0, errProductsResponse
	}

	products := []models.Product{}
	countProduct := len(productsResponse.Data)

	if countProduct == 0 {
		return 0, nil
	}

	currentItem = currentItem + pageSize
	fmt.Println("countProduct: ", countProduct)
	fmt.Println("lastItemId: ", currentItem)
	var layout = "2006-01-02T15:04:05.0000000"

	for i := 0; i < countProduct; i++ {
		product := models.Product{}
		product.UserAppTokenId = userAppToken.ID
		product.Gateway = repository.GATEWAY_KIOTVIET
		product.StoreId = store.ID
		product.ProductStatus = repository.PRODUCT_STATUS_ACTIVE
		product.Content = productsResponse.Data[i].Content
		product.Summary = productsResponse.Data[i].Summary

		createdOn, errCreatedOn := time.Parse(layout, productsResponse.Data[i].CreatedOn)
		if (errCreatedOn == nil) {
			product.CreatedOn = createdOn
		}

		product.Alias = productsResponse.Data[i].Alias
		product.ProductId = productsResponse.Data[i].ProductId

		countProductImages := len(productsResponse.Data[i].Images)
		productImages := []models.ProductImage{}
		if (countProductImages != 0) {
			for j := 0; j < countProductImages; j++ {
				productImage := models.ProductImage{}
				productImage.Src = productsResponse.Data[i].Images[j]
				productImages = append(productImages, productImage)
			}
		}
		product.Images = productImages

		productOptions := []models.Option{}
		countProductOptions := len(productsResponse.Data[i].Options)

		if (countProductOptions != 0) {
			for k := 0; k < countProductOptions; k++ {
				productOption := models.Option{}
				productOption.ProductId = productsResponse.Data[i].Options[k].ProductId
				productOption.Name = productsResponse.Data[i].Options[k].Name
				productOption.Values = append(productOption.Values, productsResponse.Data[i].Options[k].Value)
				productOptions = append(productOptions, productOption)
			}
		}
		product.Options = productOptions
		product.ProductType = productsResponse.Data[i].ProductType

		publishedOn, errPublishedOn := time.Parse(layout, productsResponse.Data[i].PublishedOn)
		if (errPublishedOn == nil) {
			product.PublishedOn = publishedOn
		}

		product.Tags = productsResponse.Data[i].Tags
		product.ProductName = productsResponse.Data[i].ProductName

		modifiedOn, errModifiedOn := time.Parse(layout, productsResponse.Data[i].ModifiedOn)
		if (errModifiedOn == nil) {
			product.ModifiedOn = modifiedOn
		}
		product.Vendor = productsResponse.Data[i].Vendor
		products = append(products, product)
	}

	if err := database.Database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_app_token_id"}, {Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"content", "summary", "alias", "images", "options", "product_type", "tags", "product_name", "modified_on", "variants", "vendor"}),
	}).Create(&products).Error; err != nil {
		return 0, err
	}

	return currentItem, nil
}
