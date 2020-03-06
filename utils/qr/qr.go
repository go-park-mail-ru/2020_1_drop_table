package qr

import (
	"github.com/skip2/go-qrcode"
	"net/url"
)

func generate(str string, qrSize int) ([]byte, error) {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return []byte(nil), err
	}

	var image []byte
	image, err = qrcode.Encode(str, qrcode.Highest, qrSize)

	if err != nil{
		return []byte(nil), err
	}
	return image, nil
}
