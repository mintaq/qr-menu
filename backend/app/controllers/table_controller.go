package controllers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func CreateTable(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	table := new(models.Table)

	if err := c.BodyParser(table); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := utils.NewValidator().Struct(table); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	store := new(models.Store)

	if tx := database.Database.Where("id = ? AND user_id = ?", table.StoreId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	if tx := database.Database.Create(table); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	tableURL := fmt.Sprintf("https://%s?table=%d", store.GetSubdomainWithSuffix(), table.ID)
	qrCodeFileName := fmt.Sprintf("%s%d.png", os.Getenv("QR_CODE_TABLE_IMAGE_PREFIX"), table.ID)
	qrCodeSrc, err := utils.CreateQRCode(store.Subdomain, tableURL, qrCodeFileName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if tx := database.Database.Model(table).Where("id = ?", table.ID).Updates(models.Table{
		QrCodeSrc: qrCodeSrc,
		TableURL:  tableURL,
	}); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  table,
	})
}

func GetTables(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	var tables []models.Table
	query := database.Database.Where("store_id = ?", c.Query("store_id"))
	tableName := c.Query("name")
	if tableName != "" {
		query = query.Where("name LIKE ?", "%"+tableName+"%")
	}

	pagination, scope := models.Paginate(models.Table{}, c, query)

	if tx := database.Database.Scopes(scope).Find(&tables); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":      false,
		"msg":        "success",
		"data":       tables,
		"pagination": pagination,
	})
}
