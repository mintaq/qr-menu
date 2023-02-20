package controllers

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func GetCollections(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	store := new(models.Store)
	if tx := database.Database.Where("id = ? AND user_id = ?", c.Query("store_id"), claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	var collections []models.Collection
	query := database.Database.Where("store_id = ?", c.Query("store_id"))
	pagination, scopes := models.Paginate(models.Collection{}, c, query)

	if tx := database.Database.Scopes(scopes).Find(&collections); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":      false,
		"msg":        "",
		"data":       collections,
		"pagination": pagination,
	})
}

func CreateCollection(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	collection := new(models.Collection)

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
			collection.StoreId = uint64(storeIdUint64)
		case "name":
			collection.Name = value[0]
		case "description":
			collection.Description = value[0]
		}
	}

	if err := utils.NewValidator().Struct(collection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	collection.Gateway = repository.GATEWAY_CUSTOM
	collection.CollectionId = utils.CreateUintId()
	collection.Alias = collection.GetNameAlias()
	store := new(models.Store)

	if tx := database.Database.First(store, "id = ? AND user_id = ?", collection.StoreId, claims.UserID); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	files := form.File["image"]
	for _, file := range files {
		if !strings.Contains(file.Header["Content-Type"][0], "image/") {
			continue
		}

		fileName := fmt.Sprintf("%s%d", os.Getenv("COLLECTION_IMAGE_PREFIX"), collection.CollectionId)
		filePathSrc, err := utils.CreateImage(file, fileName, store.Subdomain, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		collection.Image = models.CollectionImage{
			Id:  utils.CreateUintId(),
			Src: filePathSrc,
		}
	}

	if tx := database.Database.Create(collection); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "",
		"data":  collection,
	})
}
