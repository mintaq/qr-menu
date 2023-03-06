package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/cache"
)

func AddItemToCart(c *fiber.Ctx) error {
	reqHeader := new(models.Header)

	if err := c.ReqHeaderParser(reqHeader); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	if err := cache.RedisClient.Set(context.Background(), reqHeader.Cookie, "test", time.Minute*2).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", nil))
}

func GetCart(c *fiber.Ctx) error {
	reqHeader := new(models.Header)

	if err := c.ReqHeaderParser(reqHeader); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	redisCmd := cache.RedisClient.Get(context.Background(), reqHeader.Cookie)
	switch {
	case redisCmd.Err() == redis.Nil:
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse("key does not exist"))
	case redisCmd.Err() != nil:
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(redisCmd.Err().Error()))
	}

	data := map[string]interface{}{
		"redis_value": redisCmd.String(),
		"req_header":  reqHeader,
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", data))
}

func UpdateCart(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", nil))
}
