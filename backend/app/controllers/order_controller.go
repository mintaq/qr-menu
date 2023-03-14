package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/cache"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
	"gorm.io/gorm/clause"
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
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("cart is empty"))
	}

	cart.UpdateAllItemsOrderedStatus()
	order := models.Order{
		StoreId:           cart.Items[0].StoreId,
		CartToken:         cartToken,
		LineItems:         cart.GetOrderLineItems(),
		Currency:          "VND",
		FinancialStatus:   repository.FINANCIAL_STATUS_PENDING,
		Status:            repository.ORDER_STATUS_OPEN,
		FulfillmentStatus: repository.FULFILLMENT_STATUS_PARTIAL,
		Gateway:           repository.GATEWAY_CUSTOM,
	}

	if err := database.Database.Model(models.Order{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "store_id"}, {Name: "cart_token"}},
		DoUpdates: clause.AssignmentColumns(order.GetColumnsUpdateOnConflict()),
	}).Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	if err := cache.SetCartData(tableKey, cart, utils.GetRedisCartDuration()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(cart))
}

func CheckoutOrder(c *fiber.Ctx) error {
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

	// Get order from database
	order := new(models.Order)
	if err := database.Database.Where("store_id = ? AND cart_token = ?", cart.Items[0].StoreId, cartToken).First(order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(order))
}

func PayOrder(c *fiber.Ctx) error {
	// Get cart token from cookies
	cartToken := c.Cookies("cart_token")
	if cartToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("missing cart token"))
	}

	// Get cart data from cache
	tableKey, err := utils.GetHashTableKey(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}
	cart, err := cache.GetCartData(tableKey, cartToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	// Check if all items in cart are ordered
	if !cart.IsAllItemOrdered() {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("cart has unordered items"))
	}

	payOrderReq := new(models.PayOrderReqBody)
	if err := c.BodyParser(payOrderReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
	}
	if err := utils.NewValidator().Struct(payOrderReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Get order from database
	order := new(models.Order)
	if err := database.Database.Where("store_id = ? AND cart_token = ?", cart.Items[0].StoreId, cartToken).First(order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
	}

	// Check if order is already fulfilled
	if order.FulfillmentStatus == repository.FULFILLMENT_STATUS_FULFILLED {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("order is already fulfilled"))
	}

	// Update order with payment method and fulfillment status
	order.UpdatePaymentGatewayNames(payOrderReq.PaymentMethod)
	if payOrderReq.PaymentMethod == repository.PAYMENT_METHOD_CASH {
		order.FulfillmentStatus = repository.FULFILLMENT_STATUS_FULFILLED
		order.FinancialStatus = repository.FINANCIAL_STATUS_PAID
	}

	// Save updated order to database
	if err := database.Database.Model(order).Select("payment_method", "fulfillment_status", "financial_status").Updates(order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(order))
}
