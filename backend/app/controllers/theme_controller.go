package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/app/models"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/utils"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/platform/database"
)

func CreateTheme(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	theme := new(models.Theme)

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// => *multipart.Form
	if storeId := form.Value["store_id"]; len(storeId) > 0 {
		// Get key value:
		storeIdUint64, _ := strconv.Atoi(storeId[0])
		theme.StoreId = uint64(storeIdUint64)
	}

	if colors := form.Value["colors"]; len(colors) > 0 {
		// Get key value:
		themeColors := new(models.ThemeColors)
		if err := json.Unmarshal([]byte(colors[0]), themeColors); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		theme.Colors = *themeColors
	}

	if err := utils.NewValidator().Struct(theme); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	store := new(models.Store)

	if tx := database.Database.Where("id = ? AND user_id = ?", theme.StoreId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	if tx := database.Database.Create(theme); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	hostURL, _ := utils.ConnectionURLBuilder(repository.STATIC_PUBLIC_URL)
	staticPublicPath, err := utils.GetStaticPublicPathByStore(store.Subdomain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all files from "file" key:
	files := form.File["file"]

	// Loop through files:
	for _, file := range files {
		if !strings.Contains(file.Header["Content-Type"][0], "image/") {
			log.Println("test")
			continue
		}

		contentType := strings.Split(file.Header["Content-Type"][0], "/")
		imageType := contentType[1]
		fileName := fmt.Sprintf("%s%d.%s", os.Getenv("THEME_COVER_IMAGE_PREFIX"), theme.ID, imageType)
		filePath := fmt.Sprintf("%s/%s", staticPublicPath, fileName)
		filePathSrc := fmt.Sprintf("%s/stores/%s/%s", hostURL, store.Subdomain, fileName)

		// Save the files to disk:
		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		if tx := database.Database.Model(theme).Where("id = ?", theme.ID).Update("cover_image", filePathSrc); tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   tx.Error.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  theme,
	})
}

func GetThemes(c *fiber.Ctx) error {
	_, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	storeId := c.Query("store_id")
	if storeId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "store_id is missing from params",
		})
	}

	themes := new([]models.Theme)

	if tx := database.Database.Find(themes, "store_id = ?", storeId); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error": false,
		"msg":   "success",
		"data":  themes,
	})
}
