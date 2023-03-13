package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
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

// Product struct to describe product object.
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

type ProductFormData struct {
	Product
	CollectionId uint64 `json:"collection_id"`
	MenuId       uint64 `json:"menu_id" validate:"required"`
}

type ProductDBForm interface {
	GetProduct() *Product
	GetProductNameAlias() string
}

func (p *ProductFormData) GetProduct() *Product {
	test := utils.CreateUintId()
	fmt.Println(test)
	return &p.Product
}

func (p *ProductFormData) GetProductNameAlias() string {
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

func (pfd *ProductFormData) ExtractDataFromFile(ctx *fiber.Ctx, db *gorm.DB, claims *utils.TokenMetadata, ignoreKeys []string) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	if !slices.Contains(ignoreKeys, "gateway") {
		pfd.Gateway = repository.GATEWAY_CUSTOM
	}

	if pfd.StoreId <= 0 && !slices.Contains(ignoreKeys, "product_id") {
		pfd.ProductId = utils.CreateUintId()
	}

	keysMap := map[string]string{
		"store_id":      "StoreId",
		"name":          "ProductName",
		"content":       "Content",
		"price":         "Price",
		"product_type":  "ProductType",
		"collection_id": "CollectionId",
		"is_charge_tax": "IsChargeTax",
		"menu_id":       "MenuId",
	}

	for key, value := range form.Value {
		if len(value) == 0 || slices.Contains(ignoreKeys, key) {
			continue
		}
		if fieldName, ok := keysMap[key]; ok {
			field := reflect.ValueOf(pfd).Elem().FieldByName(fieldName)
			if field.IsValid() && field.CanSet() {
				switch field.Kind() {
				case reflect.Uint64:
					val, err := strconv.ParseUint(value[0], 10, 64)
					if err != nil {
						return err
					}
					field.SetUint(val)
				case reflect.Float64:
					val, err := strconv.ParseFloat(value[0], 64)
					if err != nil {
						return err
					}
					field.SetFloat(val)
				case reflect.String:
					field.SetString(value[0])
				case reflect.Bool:
					val, err := strconv.Atoi(value[0])
					if err != nil {
						return err
					}
					field.SetBool(val == 1)
				default:
					continue
				}
			}
		}
	}

	store := new(Store)

	if tx := db.First(store, "id = ? AND user_id = ?", pfd.StoreId, claims.UserID); tx.Error != nil {
		return tx.Error
	}

	for _, file := range form.File["image"] {
		if !strings.Contains(file.Header["Content-Type"][0], "image/") {
			continue
		}

		fileName := fmt.Sprintf("%s%d", os.Getenv("PRODUCT_IMAGE_PREFIX"), pfd.ProductId)
		filePathSrc, err := utils.CreateImage(file, fileName, store.GetSubdomainWithSuffix(), ctx)
		if err != nil {
			return err
		}

		image := &ProductImage{
			Id:        utils.CreateUintId(),
			Src:       filePathSrc,
			ProductId: pfd.ProductId,
		}

		pfd.Images = append(pfd.Images, *image)
	}

	return nil
}

func (pfd *ProductFormData) UpdateImagesFromFile(ctx *fiber.Ctx, db *gorm.DB, claims *utils.TokenMetadata) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	store := new(Store)

	if tx := db.First(store, "id = ? AND user_id = ?", pfd.StoreId, claims.UserID); tx.Error != nil {
		return tx.Error
	}

	if pfd.ProductId <= 0 {
		return errors.New("invalid product id")
	}

	files := form.File["image"]
	for _, file := range files {
		if !strings.Contains(file.Header["Content-Type"][0], "image/") {
			continue
		}

		fileName := fmt.Sprintf("%s%d", os.Getenv("PRODUCT_IMAGE_PREFIX"), pfd.ProductId)
		filePathSrc, err := utils.CreateImage(file, fileName, store.Subdomain, ctx)
		if err != nil {
			return err
		}

		image := &ProductImage{
			Id:        utils.CreateUintId(),
			Src:       filePathSrc,
			ProductId: pfd.ProductId,
		}

		pfd.Images = append(pfd.Images, *image)
	}

	return nil
}

func (p *Product) UpdateDataFromProductForm(pf *Product) {
	p.ProductName = pf.ProductName
	p.Alias = pf.Alias
	p.Content = pf.Content
	p.Price = pf.Price
	p.ProductType = pf.ProductType
	p.IsChargeTax = pf.IsChargeTax
	p.Images = pf.Images
}
