package qr

import (
	"2020_1_drop_table/configs"
	"github.com/skip2/go-qrcode"
	"net/url"
	"os"
)

func Generate(str string, qrSize int) ([]byte, error) {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return nil, err
	}

	var image []byte
	image, err = qrcode.Encode(str, qrcode.Highest, qrSize)

	if err != nil {
		return []byte(nil), err
	}
	return image, nil
}

func GenerateToFile(url string, uuid string) (string, error) {
	code, err := Generate(url, 256)
	if err != nil {
		return "", err
	}
	directoryPath := configs.MediaFolder + "/qr"
	os.MkdirAll(directoryPath, os.ModePerm)
	extension := ".png"
	finalPath := directoryPath + "/" + uuid + extension
	file, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	file.Write(code)
	defer file.Close()
	return finalPath, nil
}
