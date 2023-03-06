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
	cartDuration, _ := strconv.Atoi(os.Getenv("REDIS_MAX_CART_DURATION_HOURS"))
	userToken := c.Cookies("user_token")
	storeId := c.Cookies("store_id")
	tableId := c.Cookies("table_id")
	if userToken == "" || storeId == "" || tableId == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("invalid header"))
	}

	product := new(models.Product)
	reqBody := new(ReqBody)
	tableKey := fmt.Sprintf("%s:%s", storeId, tableId)
	cart := models.Cart{
		UserToken: userToken,
	}

	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	if err := utils.NewValidator().Struct(reqBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	redisCmd := cache.RedisClient.Get(context.Background(), tableKey)

	if redisCmd.Err() != redis.Nil || redisCmd.Err() != nil {
		_ = json.Unmarshal([]byte(redisCmd.Val()), &cart)
	}

	// If has product in cache -> increase quantity of product
	if cart.HasProduct(reqBody.ProductId) {
		cart.UpdateCartByProductId(reqBody.ProductId, reqBody.Quantity)
	} else {
		// If added product is not in cache -> get data from DB then create new cache
		if err := database.Database.First(product, "product_id = ? AND store_id = ?", reqBody.ProductId, storeId).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
		}

		cart.Items = append(cart.Items, models.CartItem{
			Product:  *product,
			Quantity: reqBody.Quantity,
		})
	}

	cart.UpdateCountableFields()
	dataStr, _ := json.Marshal(cart)
	expireTime := time.Duration(cartDuration) * time.Hour

	if err := cache.RedisClient.Set(context.Background(), tableKey, dataStr, expireTime).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", nil))
}

func GetCart(c *fiber.Ctx) error {
	userToken := c.Cookies("user_token")
	if userToken == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("invalid header"))
	}
	tableKey, err := utils.GetHashTableKey(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	redisCmd := cache.RedisClient.Get(context.Background(), tableKey)
	switch {
	case redisCmd.Err() == redis.Nil:
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("key does not exist"))
	case redisCmd.Err() != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(redisCmd.Err().Error()))
	}

	val := redisCmd.Val()
	cartData := models.Cart{}
	_ = json.Unmarshal([]byte(val), &cartData)
	if cartData.UserToken != userToken {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponse(false, "success", nil))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", cartData))
}

func UpdateCart(c *fiber.Ctx) error {
	cartDuration, _ := strconv.Atoi(os.Getenv("REDIS_MAX_CART_DURATION_HOURS"))
	userToken := c.Cookies("user_token")
	if userToken == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("invalid header"))
	}
	tableKey, err := utils.GetHashTableKey(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	reqBody := new(ReqBody)

	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	if err := utils.NewValidator().Struct(reqBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	redisCmd := cache.RedisClient.Get(context.Background(), tableKey)
	switch {
	case redisCmd.Err() == redis.Nil:
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("key does not exist"))
	case redisCmd.Err() != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(redisCmd.Err().Error()))
	}

	cart := models.Cart{}
	_ = json.Unmarshal([]byte(redisCmd.Val()), &cart)

	if cart.UserToken != userToken {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewResponse(false, "success", nil))
	}

	cart.UpdateCartByIndex(reqBody.Index, reqBody.Quantity).UpdateCountableFields()
	dataStr, _ := json.Marshal(cart)
	expireTime := time.Duration(cartDuration) * time.Hour

	if err := cache.RedisClient.Set(context.Background(), tableKey, dataStr, expireTime).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", cart))
}
