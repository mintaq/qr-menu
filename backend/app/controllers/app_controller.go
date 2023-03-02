package controllers

import (
	"regexp"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func CreateApp(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	newApp := new(models.App)

	if err := c.BodyParser(newApp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	if newApp.Gateway != repository.GATEWAY_KIOTVIET && newApp.Gateway != repository.GATEWAY_SAPO {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid gateway",
		})
	}

	if newApp.Gateway == repository.GATEWAY_KIOTVIET {
		if match, _ := regexp.MatchString(repository.REGEX_KIOTVIET_DOMAIN, newApp.AppName); !match {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   "invalid app name",
			})
		}
	}

	if err := utils.NewValidator().Struct(newApp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if tx := database.Database.Create(newApp); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	userAppToken := new(models.UserAppToken)
	userAppToken.UserId = uint64(claims.UserID)
	userAppToken.AppId = newApp.ID
	userAppToken.StoreDomain = newApp.AppName

	if tx := database.Database.Create(userAppToken); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   newApp,
	})
}
