package qrcode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"strings"

	"github.com/aiyaruch1320/go-qr-code/assets"
	"github.com/divan/qrlogo"
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

		logo := c.QueryParam("logo")
		if logo == "true" {
			base64Image, err := GenerateQRCodeWithLogo(qrCode.Content)
			if err != nil {
				return err
			}
			return c.JSON(http.StatusOK, map[string]string{
				"image": base64Image,
			})
		}

		base64Image, err := GenerateQRCode(qrCode.Content)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, base64Image)
	}
}

func GenerateQRCode(content string) (string, error) {
	var base64Image string
	code, err := qrcode.Encode(content, qrcode.Medium, 2048)
	if err != nil {
		return "", err
	}
	base64Image = base64.StdEncoding.EncodeToString(code)
	base64Image = "data:image/png;base64," + base64Image
	return base64Image, nil
}

func GenerateQRCodeWithLogo(content string) (string, error) {
	var base64Image string
	b64logo := assets.GetLogoBase64()

	idx := strings.Index(b64logo, ";base64,")
	if idx < 0 {
		return base64Image, fmt.Errorf("invalid logo base64")
	}
	unbase, err := base64.StdEncoding.DecodeString(b64logo[idx+8:])
	if err != nil {
		return base64Image, err
	}
	r := bytes.NewReader(unbase)
	img, err := png.Decode(r)
	if err != nil {
		return base64Image, err
	}
	code, err := qrlogo.Encode(string(content), img, 2048)
	if err != nil {
		log.Fatal(err)
	}
	base64Image = base64.StdEncoding.EncodeToString(code.Bytes())
	base64Image = "data:image/png;base64," + base64Image
	return base64Image, nil
}
