package controllers

import "github.com/gofiber/fiber/v2"

func CreateQrCode(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": true,
		"msg":   "",
	})
}
