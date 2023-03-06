package utils

import (
	"errors"
	"fmt"
	"image/color"
	"mime/multipart"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
)

func CreateQRCode(store, url, fileName string) (string, error) {
	staticPublicPath, err := GetOrCreateStaticPublicFolderPathByStore(store)
	if err != nil {
		return "", err
	}

	hostURL, _ := ConnectionURLBuilder(repository.STATIC_PUBLIC_URL)
	file := fmt.Sprintf("%s/%s", staticPublicPath, fileName)
	filePathSrc := fmt.Sprintf("%s/stores/%s/%s", hostURL, store, fileName)

	err = qrcode.WriteColorFile(url, qrcode.Highest, 256, color.White, color.Black, file)
	if err != nil {
		return "", err
	}

	return filePathSrc, nil
}

func CreateImage(file *multipart.FileHeader, fileName, storeSubdomain string, c *fiber.Ctx) (string, error) {
	if !strings.Contains(file.Header["Content-Type"][0], "image/") {
		return "", errors.New("file is not image type")
	}

	hostURL, _ := ConnectionURLBuilder(repository.STATIC_PUBLIC_URL)
	staticPublicPath, err := GetOrCreateStaticPublicFolderPathByStore(storeSubdomain)
	if err != nil {
		return "", err
	}

	contentType := strings.Split(file.Header["Content-Type"][0], "/")
	imageType := contentType[1]
	fileName = fmt.Sprintf("%s.%s", fileName, imageType)
	filePath := fmt.Sprintf("%s/%s", staticPublicPath, fileName)
	filePathSrc := fmt.Sprintf("%s/stores/%s/%s", hostURL, storeSubdomain, fileName)
	if err := c.SaveFile(file, filePath); err != nil {
		return "", err
	}

	return filePathSrc, nil
}

func CreateUintId() uint64 {
	return uint64(time.Now().UnixMilli())
}

func GetHashTableKey(c *fiber.Ctx) (string, error) {
	storeId := c.Cookies("store_id")
	tableId := c.Cookies("table_id")
	if storeId == "" || tableId == "" {
		return "", errors.New("get hash table key fail")
	}

	return fmt.Sprintf("%s:%s", storeId, tableId), nil
}
