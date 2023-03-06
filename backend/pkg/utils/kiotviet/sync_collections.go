package kiotviet

import (
	"fmt"
	"log"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"gorm.io/gorm/clause"
)

func SyncCollections(userId uint64, storeId uint64, pageSize int, currentItem int) (int, error) {
	log.Println("SyncCollections: Processing...")

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

	var collectionsResponse CollectionsResponse

	collectionsResponse, errCollectionsResponse := CollectionList(userId, pageSize, currentItem)
	if (errCollectionsResponse != nil) {
		return 0, errCollectionsResponse
	}

	collections := []models.Collection{}
	countCollection := len(collectionsResponse.Data)

	if countCollection == 0 {
		return 0, nil
	}

	currentItem = currentItem + pageSize
	fmt.Println("countCollection: ", countCollection)
	fmt.Println("lastItemId: ", currentItem)

	for i := 0; i < countCollection; i++ {
		collection := models.Collection{}
		collection.UserAppTokenId = userAppToken.ID
		collection.Gateway = repository.GATEWAY_KIOTVIET
		collection.StoreId = store.ID
		collection.CollectionId = collectionsResponse.Data[i].CollectionId
		collection.Description = collectionsResponse.Data[i].Description
		collection.Alias = collectionsResponse.Data[i].Alias
		collection.Name = collectionsResponse.Data[i].Name

		collections = append(collections, collection)
	}

	if err := database.Database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "store_id"}, {Name: "collection_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"description", "alias", "name", "image"}),
	}).Create(&collections).Error; err != nil {
		return 0, err
	}

	return currentItem, nil
}
