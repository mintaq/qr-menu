package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils/kiotviet"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

// CreateKiotvietUser method to create Kiotviet user
// @Description Create Kiotviet user
// @Summary create Kiotviet user
// @Tags KiotViet
// @Accept json
// @Produce json
// @Success 200 {string} url
// @Router /v1/kiotviet/create-user [post]
func CreateKiotvietUser(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	newApp := &models.App{}

	if err := c.BodyParser(newApp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(newApp); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	newApp.Gateway = "kiotviet"

	insertedApp := database.Database.Create(newApp)
	if insertedApp.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   insertedApp.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  newApp,
	})
}

func SyncKiotvietProducts(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	result, err := kiotviet.ConnectToken(uint64(claims.UserID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  result,
	})
}
