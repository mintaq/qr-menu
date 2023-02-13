package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func CreateProduct(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	createProductBody := new(models.CreateProductBody)

	if err := c.BodyParser(createProductBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if err := utils.NewValidator().Struct(createProductBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	createProductBody.ProductId = uint64(time.Now().Unix())
	createProductBody.Alias = createProductBody.GenProductNameAlias()
	createProductBody.Gateway = repository.GATEWAY_CUSTOM
	product := createProductBody.GetProduct()

	if tx := database.Database.Create(product); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	collect := &models.Collect{
		CollectionId: createProductBody.CollectionId,
		ProductId:    createProductBody.ProductId,
		StoreId:      createProductBody.StoreId,
	}

	if tx := database.Database.Create(collect); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  product,
	})
}
