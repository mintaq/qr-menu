package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type Collection struct {
	BasicModel
	SapoCollectionResp
	CollectionId   uint64    `json:"collection_id"`
	StoreId        uint64    `json:"store_id" validate:"required"`
	UserAppTokenId uint64    `json:"user_app_token_id" gorm:"default:null"`
	Gateway        string    `json:"gateway"`
	IsFeatured     int       `json:"is_featured" gorm:"default:0"`
	Products       []Product `json:"products" gorm:"-"`
}

type CollectionWithProducts struct {
	Collection `json:"collection"`
	Products   []Product `json:"products" gorm:"-"`
}

type SapoCollectionResp struct {
	CollectionId uint64          `json:"id"`
	Description  string          `json:"description" gorm:"default:null"`
	Alias        string          `json:"alias" gorm:"default:null"`
	Name         string          `json:"name" validate:"required"`
	Image        CollectionImage `json:"image" gorm:"default:null"`
}

type CollectionImage struct {
	Id        uint64    `json:"id"`
	Src       string    `json:"src"`
	CreatedOn time.Time `json:"created_on" gorm:"default:null"`
}

func (sla *CollectionImage) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CollectionImage) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

func (p *Collection) GetNameAlias() string {
	lowerCaseName := (strings.ToLower((strings.Trim(p.Name, " "))))
	return strings.ReplaceAll(lowerCaseName, " ", "_")
}

type CollectionQueryParams struct {
	PaginationQueryParams
	Name       string   `query:"name"`
	StoreId    uint64   `query:"store_id" validate:"required"`
	Includes   []string `query:"includes"`
	IsFeatured int      `query:"is_featured"`
}

func (c *Collection) GetProducts(db *gorm.DB) ([]Product, error) {
	products := []Product{}

	if tx := db.Model(Product{}).Joins("left join collects on collects.product_id = products.product_id and collects.collection_id = ?", c.CollectionId).Where("products.store_id = ? AND collects.collection_id = ?", c.StoreId, c.CollectionId).Find(&products); tx.Error != nil {
		return nil, tx.Error
	}

	c.Products = products

	return products, nil
}

func (c *Collection) ExtractDataFromFile(ctx *fiber.Ctx, db *gorm.DB, claims *utils.TokenMetadata, excepts []string) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	for key, value := range form.Value {
		if len(value) == 0 || slices.Contains(excepts, key) {
			continue
		}

		switch key {
		case "store_id":
			storeIdUint64, _ := strconv.Atoi(value[0])
			c.StoreId = uint64(storeIdUint64)
		case "name":
			c.Name = value[0]
			c.GetNameAlias()
		case "description":
			c.Description = value[0]
		}
	}

	if c.StoreId <= 0 {
		return errors.New("store id is invalid")
	}

	store := new(Store)

	if tx := db.First(store, "id = ? AND user_id = ?", c.StoreId, claims.UserID); tx.Error != nil {
		return tx.Error
	}

	files := form.File["image"]
	for _, file := range files {
		if !strings.Contains(file.Header["Content-Type"][0], "image/") {
			continue
		}

		fileName := fmt.Sprintf("%s%d", os.Getenv("COLLECTION_IMAGE_PREFIX"), c.CollectionId)
		filePathSrc, err := utils.CreateImage(file, fileName, store.Subdomain, ctx)
		if err != nil {
			return err
		}
		c.Image = CollectionImage{
			Id:  utils.CreateUintId(),
			Src: filePathSrc,
		}
	}

	return nil
}
