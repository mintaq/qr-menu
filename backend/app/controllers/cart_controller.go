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

func AddItemToCart(c *fiber.Ctx) error {
	cartDuration, _ := strconv.Atoi(os.Getenv("REDIS_MAX_CART_DURATION_HOURS"))
	userToken := c.Cookies("user_token")
	storeId := c.Cookies("store_id")
	tableId := c.Cookies("table_id")
	if userToken == "" || storeId == "" || tableId == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("invalid header"))
	}

	type ReqBody struct {
		ProductId uint64 `json:"product_id" validate:"required"`
		Quantity  int    `json:"quantity" validate:"required"`
	}

	product := new(models.Product)
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

	if err := database.Database.First(product, "product_id = ? AND store_id = ?", reqBody.ProductId, storeId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	cartItems := []models.CartItem{}
	cartItems = append(cartItems, models.CartItem{
		Product:  *product,
		Quantity: reqBody.Quantity,
	})
	cartData := models.Cart{
		UserToken: userToken,
		Items:     cartItems,
	}
	cartData.CountTotalItems().CountTotalPrice()

	tableKey := fmt.Sprintf("%s:%s", storeId, tableId)
	dataStr, _ := json.Marshal(cartData)
	expireTime := time.Duration(cartDuration) * time.Hour

	if err := cache.RedisClient.Set(context.Background(), tableKey, dataStr, expireTime).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", nil))
}

func GetCart(c *fiber.Ctx) error {
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
	var res interface{}
	_ = json.Unmarshal([]byte(val), &res)

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", res))
}

func UpdateCart(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", nil))
}
