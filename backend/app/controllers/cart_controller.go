package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/cache"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

type ReqBody struct {
	ProductId uint64 `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity"`
	Index     int    `json:"index"`
}

func AddItemToCart(c *fiber.Ctx) error {
	// Get cart duration from environment variable
	cartDuration, _ := strconv.Atoi(os.Getenv("REDIS_MAX_CART_DURATION_HOURS"))

	// Get required cookies
	cartToken, storeId, tableId := c.Cookies("cart_token"), c.Cookies("store_id"), c.Cookies("table_id")
	if cartToken == "" || storeId == "" || tableId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("invalid header"))
	}

	// Parse request body
	reqBody := new(ReqBody)
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
	tableKey := fmt.Sprintf("%s:%s", storeId, tableId)
	cart := models.Cart{CartToken: cartToken}
	redisCmd := cache.RedisClient.Get(context.Background(), tableKey)
	if err := redisCmd.Err(); err != nil && err != redis.Nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	} else if err := json.Unmarshal([]byte(redisCmd.Val()), &cart); err != nil {
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
	dataStr, _ := json.Marshal(cart)
	expireTime := time.Duration(cartDuration) * time.Hour
	if err := cache.RedisClient.Set(context.Background(), tableKey, dataStr, expireTime).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", nil))
}

func GetCart(c *fiber.Ctx) error {
	cartToken := c.Cookies("cart_token")
	if cartToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("missing user token"))
	}

	tableKey, err := utils.GetHashTableKey(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	cart, err := cache.GetCartData(tableKey)
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
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("missing user token"))
	}

	tableKey, err := utils.GetHashTableKey(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	reqBody := new(ReqBody)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("invalid request body"))
	}

	if err := utils.NewValidator().Struct(reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	cart, err := cache.GetCartData(tableKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	if cart.CartToken != cartToken {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse("user token does not match cart data"))
	}

	cart.UpdateCartByIndex(reqBody.Index, reqBody.Quantity).UpdateCountableFields()
	expireDuration, err := time.ParseDuration(os.Getenv("REDIS_MAX_CART_DURATION_HOURS") + "h")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	if err := cache.SetCartData(tableKey, cart, expireDuration); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", cart))
}
