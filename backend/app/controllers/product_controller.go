package controllers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func CreateProduct(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	createProductBody := new(models.CreateProductBody)
	store := new(models.Store)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	for key, value := range form.Value {
		if len(value) == 0 {
			continue
		}
		switch key {
		case "store_id":
			storeIdUint64, _ := strconv.Atoi(value[0])
			createProductBody.StoreId = uint64(storeIdUint64)
		case "name":
			createProductBody.ProductName = value[0]
		case "content":
			createProductBody.Content = value[0]
		case "price":
			createProductBody.Price, _ = strconv.ParseFloat(value[0], 64)
		case "product_type":
			createProductBody.ProductType = value[0]
		case "collection_id":
			collectionId, _ := strconv.Atoi(value[0])
			createProductBody.CollectionId = uint64(collectionId)
		}
	}
	createProductBody.Gateway = repository.GATEWAY_CUSTOM
	createProductBody.ProductId = uint64(time.Now().Unix())
	createProductBody.Alias = createProductBody.GetProductNameAlias()

	if tx := database.Database.Where("id = ? AND user_id = ?", createProductBody.StoreId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	files := form.File["image"]
	for _, file := range files {
		if !strings.Contains(file.Header["Content-Type"][0], "image/") {
			continue
		}

		fileName := fmt.Sprintf("%s%d", os.Getenv("PRODUCT_IMAGE_PREFIX"), createProductBody.ProductId)
		filePathSrc, err := utils.CreateImage(file, fileName, store.Subdomain, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		image := &models.Image{
			Id:        uint64(time.Now().Unix()),
			Src:       filePathSrc,
			ProductId: createProductBody.ProductId,
		}

		createProductBody.Images = append(createProductBody.Images, *image)
	}

	if err := utils.NewValidator().Struct(createProductBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	product := createProductBody.GetProduct()

	if tx := database.Database.Create(product); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	if createProductBody.CollectionId != 0 {
		collect := &models.Collect{
			CollectionId: createProductBody.CollectionId,
			ProductId:    createProductBody.ProductId,
			StoreId:      createProductBody.StoreId,
		}

		if tx := database.Database.Create(collect); tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   tx.Error.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  product,
	})
}

func GetProducts(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	var products []models.Product
	query := database.Database.Where("store_id = ?", c.Query("store_id"))
	pagination, scopes := models.Paginate(models.Product{}, c, query)

	if tx := database.Database.Scopes(scopes).Find(&products); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":      false,
		"msg":        "success",
		"data":       products,
		"pagination": pagination,
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	store := new(models.Store)
	productId := c.Params("product_id")
	storeId := c.Query("store_id")

	if tx := database.Database.Where("id = ? AND user_id = ?", storeId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	if tx := database.Database.Where("store_id = ? AND product_id = ?", storeId, productId).Delete(models.Collect{}); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	tx := database.Database.Where("product_id = ? AND store_id = ?", productId, storeId).Delete(models.Product{})
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":        false,
		"msg":          "success",
		"row_affected": tx.RowsAffected,
	})
}
