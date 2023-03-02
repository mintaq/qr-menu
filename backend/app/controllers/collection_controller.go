package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"golang.org/x/exp/slices"
)

func GetCollections(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	queryParams := new(models.CollectionQueryParams)

	if err := c.QueryParser(queryParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	store := new(models.Store)
	if tx := database.Database.Where("id = ? AND user_id = ?", queryParams.StoreId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	isFeatured, _ := strconv.Atoi(c.Query("is_featured", "0"))

	query := database.Database.Where("store_id = ? AND is_featured = ?", store.ID, isFeatured)
	pagination, scope := models.Paginate(models.Collection{}, c, query)
	collections := []models.Collection{}

	if tx := database.Database.Model(models.Collection{}).Scopes(scope).Find(&collections); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	if !slices.Contains(queryParams.Includes, "products") {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error":      false,
			"msg":        "success",
			"data":       collections,
			"pagination": pagination,
		})
	}

	for index := range collections {
		collectionPointer := &collections[index]
		_, err := collectionPointer.GetProducts(database.Database)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":      false,
		"msg":        "success",
		"data":       collections,
		"pagination": pagination,
	})
}

func GetFeaturedCollection(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	queryParams := new(models.CollectionQueryParams)

	if err := c.QueryParser(queryParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	store := new(models.Store)
	if tx := database.Database.Where("id = ? AND user_id = ?", queryParams.StoreId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	collection := new(models.Collection)
	if tx := database.Database.First(collection, "store_id = ? AND is_featured = 1", store.ID); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	if !slices.Contains(queryParams.Includes, "products") {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"msg":   "success",
			"data":  collection,
		})
	}

	_, _ = collection.GetProducts(database.Database)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  collection,
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
	if err := collection.ExtractDataFromFile(c, database.Database, claims, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	collection.Gateway = repository.GATEWAY_CUSTOM
	collection.CollectionId = utils.CreateUintId()
	collection.Alias = collection.GetNameAlias()

	if err := utils.NewValidator().Struct(collection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
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
