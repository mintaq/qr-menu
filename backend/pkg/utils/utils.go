package utils

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

type Utils struct {
}

type Storage interface {
	CreateImage(file *multipart.FileHeader, fileName, storeSubdomain string, c *fiber.Ctx) (string, error)
}
