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

	for key, value := range form.Value {
		if len(value) == 0 {
			continue
		}
		switch key {
		case "store_id":
			storeIdUint64, _ := strconv.Atoi(value[0])
			theme.StoreId = uint64(storeIdUint64)
		case "colors":
			themeColors := new(models.ThemeColors)
			if err := json.Unmarshal([]byte(value[0]), themeColors); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"msg":   err.Error(),
				})
			}
			theme.Colors = *themeColors
		case "role":
			theme.Role = value[0]
		}
	}

	if err := utils.NewValidator().Struct(theme); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if theme.Role == repository.THEME_ROLE_MAIN {
		var countMainTheme int64
		if tx := database.Database.Model(models.Theme{}).Where("store_id = ? AND role = ?", theme.StoreId, "main").Count(&countMainTheme); tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   tx.Error.Error(),
			})
		}

		if countMainTheme >= 1 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   "store can only have one main theme at a time",
			})
		}
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

	// Loop through files:
	for _, file := range form.File["file"] {
		if !strings.Contains(file.Header["Content-Type"][0], "image/") {
			continue
		}

		fileName := fmt.Sprintf("%s%d", os.Getenv("THEME_COVER_IMAGE_PREFIX"), theme.ID)
		filePathSrc, err := utils.CreateImage(file, fileName, store.Subdomain, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		log.Println(filePathSrc)

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

func UpdateTheme(c *fiber.Ctx) error {
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	store := new(models.Store)
	theme := new(models.Theme)
	themeId, _ := strconv.Atoi(c.Params("id"))

	if tx := database.Database.Where("id = ?", themeId).First(theme); tx.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	// Check user's store validation
	if tx := database.Database.Where("id = ? AND user_id = ?", theme.StoreId, claims.UserID).First(store); tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	for key, value := range form.Value {
		if len(value) == 0 {
			continue
		}
		switch key {
		case "colors":
			themeColors := new(models.ThemeColors)
			if err := json.Unmarshal([]byte(value[0]), themeColors); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": true,
					"msg":   err.Error(),
				})
			}
			theme.Colors = *themeColors
		case "role":
			theme.Role = value[0]
		}
	}

	for _, file := range form.File["file"] {
		if !strings.Contains(file.Header["Content-Type"][0], "image/") {
			continue
		}

		fileName := fmt.Sprintf("%s%d", os.Getenv("THEME_COVER_IMAGE_PREFIX"), theme.ID)
		filePathSrc, err := utils.CreateImage(file, fileName, store.Subdomain, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		theme.CoverImage = filePathSrc
	}

	if err := utils.NewValidator().Struct(theme); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	if theme.Role == repository.THEME_ROLE_MAIN {
		// Update other main themes to unpublished
		if tx := database.Database.Model(&models.Theme{}).Where("id <> ? AND store_id = ? AND role = ?", themeId, theme.StoreId, repository.THEME_ROLE_MAIN).Updates(models.Theme{
			Role: repository.THEME_ROLE_UNPUBLISHED,
		}); tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   tx.Error.Error(),
			})
		}
	}

	tx := database.Database.Model(theme).Where("id = ? ", themeId).Updates(models.Theme{
		Colors:     theme.Colors,
		CoverImage: theme.CoverImage,
		Role:       theme.Role,
	})
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   tx.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":        false,
		"msg":          "success",
		"row_affected": tx.RowsAffected,
	})
}
