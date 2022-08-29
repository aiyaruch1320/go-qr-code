package qrcode

import (
	"encoding/base64"
	"net/http"

	"github.com/labstack/echo/v4"
	qrcode "github.com/skip2/go-qrcode"
)

type QRCode struct {
	Content string `json:"content"`
}

func CreateQRCode() echo.HandlerFunc {
	return func(c echo.Context) error {
		qrCode := new(QRCode)
		if err := c.Bind(qrCode); err != nil {
			return err
		}
		png, err := GenerateQRCode(qrCode.Content)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, png)
	}
}

func GenerateQRCode(content string) (string, error) {
	var base64Image string
	code, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64Image = base64.StdEncoding.EncodeToString(code)
	base64Image = "data:image/png;base64," + base64Image
	return base64Image, nil
}
