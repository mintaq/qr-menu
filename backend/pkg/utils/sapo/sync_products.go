package sapo

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

func SyncProducts(page, limit int, store string) (int, error) {
	log.Println("SyncProducts: Processing...")

	foundStore := new(models.Store)

	if tx := database.Database.First(foundStore, "store = ?", store); tx.Error != nil {
		return 0, tx.Error
	}

	requestURI := fmt.Sprintf("https://%s/admin/products.json?page=%d&limit=%d", store, page, limit)
	log.Println(requestURI)
	req, err := http.NewRequest(http.MethodGet, requestURI, http.NoBody)
	if err != nil {
		return 0, err
	}
	req.Header.Set("X-Sapo-Access-Token", foundStore.AccessToken)
	resp, err := utils.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	type RespProducts struct {
		Products []models.SapoProductResp
	}

	respProducts := new(RespProducts)

	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(respProducts)
		if err != nil {
			return 0, err
		}
	}

	products := []models.Product{}
	countProduct := len(respProducts.Products)

	if countProduct == 0 {
		return countProduct, nil
	}

	for i := 0; i < countProduct; i++ {
		product := models.Product{}
		product.StoreId = foundStore.ID
		product.Gateway = repository.GATEWAY_SAPO
		product.SapoProductResp = respProducts.Products[i]
		product.ProductId = respProducts.Products[i].ProductId

		products = append(products, product)
	}

	if err := database.Database.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "store_id"}, {Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"content", "summary", "alias", "images", "options", "product_type", "tags", "product_name", "modified_on", "variants", "vendor"}),
	}).Create(&products).Error; err != nil {
		return 0, err
	}

	return countProduct, nil
}
