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

type SapoProductResp struct {
	Content     string       `json:"content" validate:"required"`
	Summary     string       `json:"summary" gorm:"default:null"`
	CreatedOn   time.Time    `json:"created_on" gorm:"default:null"`
	Alias       string       `json:"alias"`
	ProductId   uint64       `json:"id"`
	Images      ImageArray   `json:"images" gorm:"default:null"`
	Options     OptionArray  `json:"options" gorm:"default:null"`
	ProductType string       `json:"product_type" gorm:"default:null"`
	PublishedOn time.Time    `json:"published_on" gorm:"default:null"`
	Tags        string       `json:"tags" gorm:"default:null"`
	ProductName string       `json:"name" gorm:"column:product_name" validate:"required"`
	ModifiedOn  time.Time    `json:"modified_on" gorm:"default:null"`
	Variants    VariantArray `json:"variants"`
	Vendor      string       `json:"vendor" gorm:"default:null"`
}

// User struct to describe Product object.
type Product struct {
	BasicModel
	SapoProductResp
	ProductId      uint64    `json:"product_id"`
	Price          float64   `json:"price" validate:"required" gorm:"default:null"`
	StoreId        uint64    `json:"store_id" validate:"required"`
	UserAppTokenId uint64    `json:"user_app_token_id" gorm:"default:null"`
	IsChargeTax    int       `json:"is_charge_tax" validate:"eq=0|eq=1"`
	ProductStatus  string    `json:"product_status"`
	Gateway        string    `json:"gateway" validate:"required"`
	Status         string    `json:"status" gorm:"default:'active'"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateProductBody struct {
	Product
	CollectionId uint64 `json:"collection_id"`
	MenuId       uint64 `json:"menu_id" validate:"required"`
}

type ProductDBForm interface {
	GetProduct() *Product
	GetProductNameAlias() string
}

func (p *CreateProductBody) GetProduct() *Product {
	test := utils.CreateUintId()
	fmt.Println(test)
	return &p.Product
}

func (p *CreateProductBody) GetProductNameAlias() string {
	lowerCaseName := (strings.ToLower((strings.Trim(p.ProductName, " "))))
	return strings.ReplaceAll(lowerCaseName, " ", "_")
}

func (p *Product) GetProductNameAlias() string {
	lowerCaseName := (strings.ToLower((strings.Trim(p.ProductName, " "))))
	p.Alias = strings.ReplaceAll(lowerCaseName, " ", "_")

	return p.Alias
}

type VariantArray []Variant

func (sla *VariantArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla VariantArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type ImageArray []ProductImage

func (sla *ImageArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ImageArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type OptionArray []Option

func (sla *OptionArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla OptionArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type Variant struct {
	Id                  int       `json:"id"`
	ProductId           uint64    `json:"product_id"`
	Title               string    `json:"title"`
	Price               float32   `json:"price"`
	Sku                 string    `json:"sku"`
	Position            int       `json:"position"`
	Grams               int       `json:"grams"`
	InventoryManagement string    `json:"inventory_management"`
	Option1             string    `json:"option1"`
	Option2             string    `json:"option2"`
	Option3             string    `json:"option3"`
	CreatedOn           time.Time `json:"created_on"`
	ModifiedOn          time.Time `json:"modified_on"`
	RequiresShipping    bool      `json:"requires_shipping"`
	Barcode             string    `json:"barcode"`
	InventoryQuantity   int       `json:"inventory_quantity"`
	ImageId             int       `json:"image_id"`
	Weight              float32   `json:"weight"`
	WeightUnit          string    `json:"weight_unit"`
}

type Option struct {
	Id        uint64   `json:"id"`
	ProductId uint64   `json:"product_id"`
	Name      string   `json:"name"`
	Position  int      `json:"position"`
	Values    []string `json:"values"`
}

type ProductImage struct {
	Id         uint64    `json:"id" gorm:"default:null"`
	ProductId  uint64    `json:"product_id"`
	Position   int       `json:"position" gorm:"default:null"`
	CreatedOn  time.Time `json:"created_on" gorm:"default:null"`
	ModifiedOn time.Time `json:"modified_on" gorm:"default:null"`
	Src        string    `json:"src"`
	VariantIds []int     `json:"variant_ids" gorm:"default:null"`
}

func (c *CreateProductBody) ExtractDataFromFile(ctx *fiber.Ctx, db *gorm.DB, claims *utils.TokenMetadata, excepts []string) error {
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
			c.ProductName = value[0]
			c.GetProductNameAlias()
		case "content":
			c.Content = value[0]
		case "price":
			c.Price, _ = strconv.ParseFloat(value[0], 64)
		case "product_type":
			c.ProductType = value[0]
		case "collection_id":
			collectionId, _ := strconv.Atoi(value[0])
			c.CollectionId = uint64(collectionId)
		case "is_charge_tax":
			isChargeTax, _ := strconv.Atoi(value[0])
			if isChargeTax == 0 || isChargeTax == 1 {
				c.IsChargeTax = isChargeTax
			}
		case "menu_id":
			menuId, _ := strconv.Atoi(value[0])
			c.MenuId = uint64(menuId)
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

		fileName := fmt.Sprintf("%s%d", os.Getenv("PRODUCT_IMAGE_PREFIX"), c.ProductId)
		filePathSrc, err := utils.CreateImage(file, fileName, store.Subdomain, ctx)
		if err != nil {
			return err
		}

		image := &ProductImage{
			Id:        utils.CreateUintId(),
			Src:       filePathSrc,
			ProductId: c.ProductId,
		}

		c.Images = append(c.Images, *image)
	}

	return nil
}
