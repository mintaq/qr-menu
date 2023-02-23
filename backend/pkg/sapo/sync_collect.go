package sapo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"gorm.io/gorm/clause"
)

func SyncCollect(page, limit int, sapoDomain string, storeId uint64) (int, error) {
	log.Println("SyncCollect: Processing...")
	userAppToken := new(models.UserAppToken)

	if tx := database.Database.First(userAppToken, "store_domain = ?", sapoDomain); tx.Error != nil {
		return 0, tx.Error
	}

	requestURI := fmt.Sprintf("https://%s/admin/collects.json?page=%d&limit=%d", sapoDomain, page, limit)
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

	type RespBody struct {
		Collects []models.Collect `json:"collects"`
	}

	respBody := new(RespBody)

	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(respBody)
		if err != nil {
			return 0, err
		}
	}

	collects := []models.Collect{}
	countCollect := len(respBody.Collects)

	if countCollect == 0 {
		return 0, nil
	}

	for i := 0; i < countCollect; i++ {
		collection := models.Collect{}
		collection.UserAppTokenId = userAppToken.ID
		collection.CollectionId = respBody.Collects[i].CollectionId
		collection.ProductId = respBody.Collects[i].ProductId
		collection.StoreId = storeId

		collects = append(collects, collection)
	}

	if err := database.Database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}, {Name: "collection_id"}},
		DoNothing: true,
	}).Create(&collects).Error; err != nil {
		return 0, err
	}

	return countCollect, nil
}
