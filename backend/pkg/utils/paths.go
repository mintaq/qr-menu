package utils

import (
	"fmt"
	"os"
)

func GetStaticPublicPathByStore(store string) (string, error) {
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
