package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/cache"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
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

	if cart.ItemsCount == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("can't checkout with empty cart"))
	}

	cart.UpdateAllItemsOrderedStatus()
	order := models.Order{
		StoreId:           cart.Items[0].StoreId,
		CartToken:         cartToken,
		Currency:          "VND",
		FinancialStatus:   repository.FINANCIAL_STATUS_PENDING,
		Status:            repository.ORDER_STATUS_OPEN,
		FulfillmentStatus: repository.FULFILLMENT_STATUS_NULL,
		Gateway:           repository.GATEWAY_CUSTOM,
	}

	if err := database.Database.Save(order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	if err := cache.SetCartData(tableKey, cart, utils.GetRedisCartDuration()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(cart))
}

func Checkout(c *fiber.Ctx) error {
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

	if !cart.IsAllItemOrdered() {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("cart has unordered items"))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(nil))
}
