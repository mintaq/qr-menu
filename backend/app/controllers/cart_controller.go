package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/cache"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

type CartReqBody struct {
	ProductId uint64 `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity"`
	Index     int    `json:"index"`
}

func AddItemToCart(c *fiber.Ctx) error {
	// Get required cookies
	cartToken, storeId, tableId := c.Cookies("cart_token"), c.Cookies("store_id"), c.Cookies("table_id")
	if cartToken == "" || storeId == "" || tableId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("invalid header"))
	}

	// Parse request body
	reqBody := new(CartReqBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
	}

	// Validate request body
	if err := utils.NewValidator().Struct(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Get cart from Redis cache
	tableKey, err := utils.GetHashTableKey(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	cart, err := cache.GetCartData(tableKey, cartToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	// Check if product is already in cart
	if cart.HasProduct(reqBody.ProductId) {
		cart.UpdateCartByProductId(reqBody.ProductId, reqBody.Quantity)
	} else {
		// Get product data from database
		product := new(models.Product)
		if err := database.Database.First(product, "product_id = ? AND store_id = ?", reqBody.ProductId, storeId).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
		}

		// Add new product to cart
		cart.Items = append(cart.Items, models.CartItem{
			Product:  *product,
			Quantity: reqBody.Quantity,
		})
	}

	// Update countable fields and save cart to Redis cache
	cart.UpdateCountableFields()

	if err := cache.SetCartData(tableKey, cart, utils.GetRedisCartDuration()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(cart))
}

func GetCart(c *fiber.Ctx) error {
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

	if cart.CartToken != cartToken {
		return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", cart))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", cart))
}

func UpdateCart(c *fiber.Ctx) error {
	cartToken := c.Cookies("cart_token")
	if cartToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("missing cart token"))
	}

	tableKey, err := utils.GetHashTableKey(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	type CartUpdateReqBody struct {
		Index    int `json:"index" validate:"gte=0"`
		Quantity int `json:"quantity" validate:"gte=0"`
	}

	reqBody := new(CartUpdateReqBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("invalid request body"))
	}

	if err := utils.NewValidator().Struct(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	cart, err := cache.GetCartData(tableKey, cartToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	cart.UpdateCartByIndex(reqBody.Index, reqBody.Quantity).UpdateCountableFields()

	if err := cache.SetCartData(tableKey, cart, utils.GetRedisCartDuration()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(cart))
}
