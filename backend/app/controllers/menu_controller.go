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
		"msg":   "success",
		"data":  menu,
	})
}

func GetMenus(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	var menus []models.Menu
	query := database.Database.Where("store_id = ?", c.Query("store_id"))
	menuName := c.Query("name")
	if menuName != "" {
		query = query.Where("name LIKE ?", "%"+menuName+"%")
	}

	pagination, scope := models.Paginate(models.Menu{}, c, query)

	if tx := database.Database.Scopes(scope).Find(&menus); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":      false,
		"msg":        "success",
		"data":       menus,
		"pagination": pagination,
	})
}

func DeleteMenu(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	store := new(models.Store)
	menuId := c.Params("id")
	storeId := c.Query("store_id")

	if tx := database.Database.Where("id = ? AND user_id = ?", storeId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	tx := database.Database.Where("id = ? AND store_id = ?", menuId, storeId).Delete(&models.Menu{})
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	rowsAffected := tx.RowsAffected

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":        false,
		"msg":          "success",
		"row_affected": rowsAffected,
	})
}
