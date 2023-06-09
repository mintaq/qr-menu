package utils

import (
	"fmt"
	"os"
)

// Return absolute path to store static public folder.
//
// Eg: /home/shopify/qr-menu-backend/backend/static/public/stores/example.com.vn
func GetOrCreateStaticPublicFolderPathByStore(store string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	folderPath := fmt.Sprintf("%s%s/stores/%s", dir, os.Getenv("STATIC_PUBLIC_PATH"), store)

	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return folderPath, nil
}
