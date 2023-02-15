package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
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
