package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"gorm.io/gorm"
)

// @Summary Create Collect
// @Description Create a new collect with given parameters.
// @Tags Collects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer {token}"
// @Param store_id query int true "Store ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/collects [post]
func CreateCollect(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
	}

	collect := new(models.Collect)

	if err := c.BodyParser(collect); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	if err := utils.NewValidator().Struct(collect); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if err := database.Database.Transaction(func(db *gorm.DB) error {
		store := new(models.Store)
		collection := new(models.Collection)
		product := new(models.Product)

		if err := db.First(store, "id = ? AND user_id = ?", collect.StoreId, claims.UserID).Error; err != nil {
			return err
		}

		if err := db.First(collection, "collection_id = ? AND store_id = ?", collect.CollectionId, collect.StoreId).Error; err != nil {
			return err
		}

		if err := db.First(product, "product_id = ? AND store_id = ?", collect.ProductId, collect.StoreId).Error; err != nil {
			return err
		}

		if err := db.Create(collect).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", collect))
}
