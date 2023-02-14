package kiotviet

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"gorm.io/gorm/clause"
)

func SyncCustomCollections(page, limit int, storeDomain string) (int, error) {
	log.Println("SyncCustomCollections: Processing...")
	userAppToken := new(models.UserAppToken)

	if tx := database.Database.First(userAppToken, "store_domain = ?", storeDomain); tx.Error != nil {
		return 0, tx.Error
	}

	requestURI := fmt.Sprintf("https://%s/admin/custom_collections.json?page=%d&limit=%d", storeDomain, page, limit)
	log.Println(requestURI)
	req, err := http.NewRequest(http.MethodGet, requestURI, http.NoBody)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-Sapo-Access-Token", userAppToken.AccessToken)
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	type RespCollections struct {
		CustomCollections []models.SapoCollectionResp `json:"custom_collections"`
	}

	respCollections := new(RespCollections)

	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(respCollections)
		if err != nil {
			return 0, err
		}
	}

	collections := []models.Collection{}
	countCollection := len(respCollections.CustomCollections)

	if countCollection == 0 {
		return 0, nil
	}

	for i := 0; i < countCollection; i++ {
		collection := models.Collection{}
		collection.UserAppTokenId = userAppToken.ID
		collection.Gateway = repository.GATEWAY_SAPO
		collection.SapoCollectionResp = respCollections.CustomCollections[i]
		collection.CollectionId = respCollections.CustomCollections[i].CollectionId

		collections = append(collections, collection)
	}

	if err := database.Database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_app_token_id"}, {Name: "collection_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"description", "alias", "name", "image"}),
	}).Create(&collections).Error; err != nil {
		return 0, err
	}

	return countCollection, nil
}

func SyncSmartCollections(page, limit int, storeDomain string) (int, error) {
	log.Println("SyncSmartCollections: Processing...")

	userAppToken := new(models.UserAppToken)

	if tx := database.Database.First(userAppToken, "store = ?", storeDomain); tx.Error != nil {
		return 0, tx.Error
	}

	requestURI := fmt.Sprintf("https://%s/admin/smart_collections.json?page=%d&limit=%d", storeDomain, page, limit)
	log.Println(requestURI)
	req, err := http.NewRequest(http.MethodGet, requestURI, http.NoBody)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-Sapo-Access-Token", userAppToken.AccessToken)
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	type RespCollections struct {
		SmartCollections []models.SapoCollectionResp `json:"smart_collections"`
	}

	respCollections := new(RespCollections)

	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(respCollections)
		if err != nil {
			return 0, err
		}
	}

	collections := []models.Collection{}
	countCollection := len(respCollections.SmartCollections)

	if countCollection == 0 {
		return 0, nil
	}

	for i := 0; i < countCollection; i++ {
		collection := models.Collection{}
		collection.UserAppTokenId = userAppToken.ID
		collection.Gateway = repository.GATEWAY_SAPO
		collection.SapoCollectionResp = respCollections.SmartCollections[i]
		collection.CollectionId = respCollections.SmartCollections[i].CollectionId

		collections = append(collections, collection)
	}

	if err := database.Database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_app_token_id"}, {Name: "collection_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"description", "alias", "name", "image"}),
	}).Create(&collections).Error; err != nil {
		return 0, err
	}

	return countCollection, nil
}
