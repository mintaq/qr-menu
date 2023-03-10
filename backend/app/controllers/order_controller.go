package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/cache"
)

func CreateOrder(c *fiber.Ctx) error {
	cartToken := c.Cookies("cart_token")
	if cartToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("missing cart token"))
	}

	tableKey, err := utils.GetHashTableKey(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	cart, err := cache.GetCartData(tableKey, cartToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	cart.UpdateAllItemsOrderedStatus()

	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(cart))
}
