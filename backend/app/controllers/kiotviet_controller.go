package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils/kiotviet"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/worker"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/worker/tasks"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

// CreateKiotvietUser method to create Kiotviet user
// @Description Create Kiotviet user
// @Summary create Kiotviet user
// @Tags KiotViet
// @Accept json
// @Produce json
// @Success 200 {string} url
// @Router /v1/kiotviet/create-user [post]
func CreateKiotvietUser(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	newApp := &models.App{}

	if err := c.BodyParser(newApp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(newApp); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	newApp.Gateway = "kiotviet"

	insertedApp := database.Database.Create(newApp)
	if insertedApp.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   insertedApp.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  newApp,
	})
}

type SyncKiotVietRequest struct {
	StoreId    uint64 `json:"store_id" validate:"required"`
}

func SyncKiotvietProducts(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	syncKiotVietRequest := new(SyncKiotVietRequest)

	if err := c.BodyParser(syncKiotVietRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(syncKiotVietRequest); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	task, err := tasks.NewSyncKiotVietProductsRecursiveTask(uint64(claims.UserID), uint64(syncKiotVietRequest.StoreId), 100, 0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	info, err := worker.AsynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(1*time.Minute))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  info.CompletedAt,
	})
}

func SyncKiotvietCollections(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	syncKiotVietRequest := new(SyncKiotVietRequest)

	if err := c.BodyParser(syncKiotVietRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(syncKiotVietRequest); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	kiotviet.SyncCollections(uint64(claims.UserID), uint64(syncKiotVietRequest.StoreId), 100, 0)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  "1",
	})
	// task, err := tasks.NewSyncKiotVietProductsRecursiveTask(uint64(claims.UserID), uint64(syncKiotVietRequest.StoreId), 100, 0)
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   err.Error(),
	// 	})
	// }

	// info, err := worker.AsynqClient.Enqueue(task, asynq.MaxRetry(3), asynq.Timeout(1*time.Minute))
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": true,
	// 		"msg":   err.Error(),
	// 	})
	// }

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"error": false,
	// 	"msg":   "success",
	// 	"data":  info.CompletedAt,
	// })
}
