package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
)

func CreateOrder(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(nil))
}
