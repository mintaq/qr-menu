package sapo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"gorm.io/gorm/clause"
)

func SyncCustomCollections(page, limit int, store string) error {
	foundStore := new(models.Store)

	if tx := database.Database.First(foundStore, "store = ?", store); tx.Error != nil {
		return tx.Error
	}

	requestURI := fmt.Sprintf("https://%s/admin/custom_collections.json?page=%d&limit=%d", store, page, limit)
	req, err := http.NewRequest(http.MethodGet, requestURI, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("X-Sapo-Access-Token", foundStore.AccessToken)
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type RespCollections struct {
		CustomCollections []models.SapoCollectionResp `json:"custom_collections"`
	}

	respCollections := new(RespCollections)

	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(respCollections)
		if err != nil {
			return err
		}
	}

	collections := []models.Collection{}

	for i := 0; i < len(respCollections.CustomCollections); i++ {
		collection := models.Collection{}
		collection.StoreId = foundStore.ID
		collection.Gateway = repository.GATEWAY_SAPO
		collection.SapoCollectionResp = respCollections.CustomCollections[i]
		collection.CollectionId = respCollections.CustomCollections[i].CollectionId

		collections = append(collections, collection)
	}

	if err := database.Database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "store_id"}, {Name: "collection_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"description", "alias", "name", "image"}),
	}).Create(&collections).Error; err != nil {
		return err
	}

	return nil
}
