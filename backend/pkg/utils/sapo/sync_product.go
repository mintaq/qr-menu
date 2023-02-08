package sapo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func SyncProducts(page, limit int, store string) error {
	foundStore := new(models.Store)

	if tx := database.Database.First(foundStore, "store = ?", store); tx.Error != nil {
		return tx.Error
	}

	requestURI := fmt.Sprintf("https://%s/admin/products.json", store)
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

	type RespProducts struct {
		Products []models.Product
	}

	respProducts := new(RespProducts)

	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(respProducts)
		if err != nil {
			return err
		}

	}

	return nil
}
