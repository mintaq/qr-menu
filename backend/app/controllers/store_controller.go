package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"gorm.io/gorm"
)

func CreateStore(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	store := &models.Store{}
	store.UserId = uint64(claims.UserID)

	if err := c.BodyParser(store); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := utils.NewValidator().Struct(store); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if err := database.Database.Transaction(func(db *gorm.DB) error {
		if err := db.Create(store).Error; err != nil {
			return err
		}

		collection := models.NewFeaturedCollection(store.ID)
		if err := db.Save(collection).Error; err != nil {
			return err
		}

		menu := models.NewDefaultMenu(store.ID)
		if err := db.Save(menu).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  store,
	})
}

func GetStores(c *fiber.Ctx) error {
	tokenMetaData, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	stores := []models.Store{}

	query := database.Database.Model(models.Store{}).Where("user_id = ?", tokenMetaData.UserID)

	pagination, scopes := models.Paginate(models.Store{}, c, query)

	if tx := database.Database.Model(models.Store{}).Scopes(scopes).Find(&stores); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":      false,
		"msg":        "success",
		"data":       stores,
		"pagination": pagination,
	})
}

func GetStoreById(c *fiber.Ctx) error {
	tokenMetaData, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	storeId, _ := strconv.Atoi(c.Params("id", ""))
	store := new(models.Store)

	if tx := database.Database.First(store, "id = ? AND user_id = ?", storeId, tokenMetaData.UserID); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  store,
	})
}

func GetStoreBySubdomain(c *fiber.Ctx) error {
	tokenMetaData, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	subdomain := c.Params("subdomain", "")
	store := new(models.Store)

	if tx := database.Database.First(store, "subdomain = ? AND user_id = ?", subdomain, tokenMetaData.UserID); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  store,
	})
}

func UpdateStore(c *fiber.Ctx) error {
	tokenMetaData, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	storeId, _ := strconv.Atoi(c.Params("id", ""))
	if storeId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "id is invalid",
		})
	}
	store := new(models.Store)
	store.ID = uint64(storeId)

	if err := database.Database.First(store, "id = ?", storeId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	updateData := new(models.StoreUpdatableData)
	store.UserId = uint64(tokenMetaData.UserID)

	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := utils.NewValidator().Struct(updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	store.UpdateStore(updateData)

	if tx := database.Database.Save(store); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  store,
	})
}
