package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func CreateBusiness(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	business := &models.Business{}
	business.UserId = uint64(claims.UserID)

	if err := c.BodyParser(business); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := utils.NewValidator().Struct(business); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if tx := database.Database.Create(business); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  business,
	})
}
