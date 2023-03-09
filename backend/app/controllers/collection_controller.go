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

	query := database.Database.Model(models.Collection{}).Where("store_id = ? AND is_featured = ?", store.ID, isFeatured)
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

// @Summary Retrieves the featured collection for a user's store.
// @Description This endpoint retrieves the featured collection for a user's store based on the provided query parameters. The query parameters can include an "includes" parameter to include the collection's products.
// @Tags Collections
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param store_id query int true "ID of the store to retrieve the featured collection for"
// @Param includes query []string false "Array of strings specifying what data to include in the response. Valid values are: products."
// @Success 200 {object} models.Collection "success"
// @Failure 400 {object} models.Response "error"
// @Failure 500 {object} models.Response "error"
// @Router /collections/featured [get]
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

// @Summary Creates a new collection for a user's store.
// @Description This endpoint creates a new collection for a user's store based on the data provided in a file sent in the request body. The file should be in CSV or JSON format and contain the necessary data for creating a collection. The user must be authenticated and authorized to create collections for the specified store.
// @Tags Collections
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer {token}"
// @Param file formData file true "CSV or JSON file containing the collection data"
// @Success 200 {object} models.Collection "success"
// @Failure 400 {object} models.Response "error"
// @Failure 500 {object} models.Response "error"
// @Router /collections [post]
func CreateCollection(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewResponse(true, err.Error(), nil))
	}

	collection := new(models.Collection)
	if err := collection.ExtractDataFromFile(c, database.Database, claims, nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponse(true, err.Error(), nil))
	}
	collection.Gateway = repository.GATEWAY_CUSTOM
	collection.CollectionId = utils.CreateUintId()
	collection.Alias = collection.GetNameAlias()

	if err := utils.NewValidator().Struct(collection); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponse(true, err.Error(), nil))
	}

	if tx := database.Database.Create(collection); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewResponse(true, tx.Error.Error(), nil))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", collection))
}
