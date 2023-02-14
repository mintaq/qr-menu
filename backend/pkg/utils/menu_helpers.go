package utils

import (
	"fmt"
	"image/color"

	"github.com/skip2/go-qrcode"
	"gitlab.xipat.com/omega-team3/qr-menu-backend/pkg/repository"
)

func CreateQRCode(store, url, fileName string) (string, error) {
	staticPublicPath, err := GetStaticPublicPathByStore(store)
	if err != nil {
		return "", err
	}

	hostURL, _ := ConnectionURLBuilder(repository.STATIC_PUBLIC_URL)
	file := fmt.Sprintf("%s/%s", staticPublicPath, fileName)
	filePathSrc := fmt.Sprintf("%s/stores/%s/%s", hostURL, store, fileName)

	err = qrcode.WriteColorFile(url, qrcode.Highest, 256, color.Black, color.White, file)
	if err != nil {
		return "", err
	}

	return filePathSrc, nil

}
