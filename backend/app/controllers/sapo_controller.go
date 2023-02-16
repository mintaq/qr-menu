package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils/sapo"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/worker"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/worker/tasks"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

// GetSapoAccessToken method to get Sapo access token
// @Description Get Sapo user access token
// @Summary get access token
// @Tags Sapo
// @Accept json
// @Produce json
// @Success 200 {string} token
// @Router /v1/sapo/get-token [get]
func GetSapoAccessToken(c *fiber.Ctx) error {
	code := c.Query("code")
	store := c.Query("store")
	userId := c.Query("user_id")
	if code == "" || store == "" || userId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "userId or code or store not found in query!",
		})
	}

	var app models.App
	tx := database.Database.First(&app, "gateway = ?", repository.GATEWAY_SAPO)
	if tx.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "app not found",
		})
	}

	uri := fmt.Sprintf("https://%s/admin/oauth/access_token", store)

	type Payload struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}

	payload := Payload{
		ClientId:     app.ApiKey,
		ClientSecret: app.SecretKey,
		Code:         code,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(body)) // #nosec
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	defer resp.Body.Close()

	type AccessToken struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
	}

	accessToken := new(AccessToken)

	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(accessToken)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
	}

	userAppToken := new(models.UserAppToken)
	userAppToken.UserId, _ = strconv.ParseUint(userId, 10, 64)
	userAppToken.AppId = app.ID
	userAppToken.StoreDomain = store
	userAppToken.AccessToken = accessToken.AccessToken
	if tx := database.Database.Create(userAppToken); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "success",
	})
}

// GetSapoAuthURL method to get Sapo authenticate URL
// @Description Get Sapo authenticate URL
// @Summary get authenticate URL
// @Tags Sapo
// @Accept json
// @Produce json
// @Success 200 {string} url
// @Router /v1/sapo/get-auth-url [get]
func GetSapoAuthURL(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	store := c.Query("store")
	if store == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "store not found in query!",
		})
	}

	match, _ := regexp.MatchString(repository.REGEX_SAPO_DOMAIN, store)
	if !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "store name is invalid!",
		})
	}

	url, err := sapo.GetAuthURLByStore(store, uint64(claims.UserID))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   url,
	})
}

type SyncReq struct {
	SapoDomain string `json:"sapo_domain" validate:"required"`
	StoreId    uint64 `json:"store_id" validate:"required"`
}

func SyncSapoProducts(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	syncReq := new(SyncReq)

	if err := c.BodyParser(syncReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(syncReq); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	task, err := tasks.NewSyncSapoProductsRecursiveTask(1, 1, syncReq.SapoDomain, syncReq.StoreId)
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

func SyncSapoCustomCollections(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	syncReq := new(SyncReq)

	if err := c.BodyParser(syncReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(syncReq); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	task, err := tasks.NewSyncSapoCustomCollectionsRecursiveTask(1, 1, syncReq.SapoDomain, syncReq.StoreId)
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

	syncSapoCollectRecursiveTask, err := tasks.NewSyncSapoCollectRecursiveTask(1, 1, syncReq.SapoDomain, syncReq.StoreId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	_, err = worker.AsynqClient.Enqueue(syncSapoCollectRecursiveTask, asynq.MaxRetry(3), asynq.Timeout(1*time.Minute))
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

func SyncSapoSmartCollections(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	syncReq := new(SyncReq)

	if err := c.BodyParser(syncReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	validate := utils.NewValidator()
	if err := validate.Struct(syncReq); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	task, err := tasks.NewSyncSapoSmartCollectionsRecursiveTask(1, 1, syncReq.SapoDomain, syncReq.StoreId)
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
