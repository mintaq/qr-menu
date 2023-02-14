package controllers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func CreateMenu(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	menu := new(models.Menu)

	if err := c.BodyParser(menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := utils.NewValidator().Struct(menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	store := new(models.Store)

	if tx := database.Database.Where("id = ? AND user_id = ?", menu.StoreId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	if tx := database.Database.Create(menu); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	menuURL := fmt.Sprintf("https://%s?menu=%d", store.Subdomain, menu.ID)
	qrCodeFileName := fmt.Sprintf("%s%d.png", os.Getenv("QR_CODE_MENU_IMAGE_PREFIX"), menu.ID)
	qrCodeSrc, err := utils.CreateQRCode(store.Subdomain, menuURL, qrCodeFileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if tx := database.Database.Model(menu).Where("id = ?", menu.ID).Updates(models.Menu{
		QrCodeSrc: qrCodeSrc,
		MenuURL:   menuURL,
	}); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "",
		"data":  menu,
	})
}

func ListMenus(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	menu := new(models.Menu)
	pagination := models.Pagination{
		Limit: 5,
		Page:  1,
	}
	res, err := menu.List(pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "",
		"data":  res,
	})
}
