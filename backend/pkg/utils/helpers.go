package utils

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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
	storeID := c.Cookies("store_id")
	tableID := c.Cookies("table_id")

	if storeID == "" {
		return "", errors.New("missing store ID")
	}

	if tableID == "" {
		return "", errors.New("missing table ID")
	}

	hashTableKey := fmt.Sprintf("%s:%s:%s", os.Getenv("QR_MENU_STORE_TABLE_PREFIX"), storeID, tableID)

	return hashTableKey, nil
}

func GetRedisCartDuration() time.Duration {
	expireDuration, err := time.ParseDuration(os.Getenv("REDIS_MAX_CART_DURATION_HOURS") + "h")
	if err != nil {
		log.Println(err.Error())
	}

	return expireDuration
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next(w, r)
	}
}
