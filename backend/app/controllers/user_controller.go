package controllers

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func GetUserProfile(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.NewErrorResponse(err.Error()))
	}

	user := new(models.User)

	if err := database.Database.First(user, "id = ?", claims.UserID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
	}

	user.PasswordHash = ""
	userApps := new([]models.UserAppToken)

	if err := database.Database.Find(userApps, "user_id = ?", user.ID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
	}

	type App struct {
		Id      uint64 `json:"id"`
		Gateway string `json:"gateway"`
	}
	apps := []App{}

	for _, userApp := range *userApps {
		app := models.App{}
		if err := database.Database.First(&app, "id = ?", userApp.AppId).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(models.NewErrorResponse(err.Error()))
		}
		apps = append(apps, App{Id: app.ID, Gateway: app.Gateway})
	}

	return c.Status(fiber.StatusOK).JSON(models.NewResponse(false, "success", map[string]interface{}{
		"user": user,
		"apps": apps,
	}))
}
