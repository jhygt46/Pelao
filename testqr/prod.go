package main

import (
	"fmt"
	"image/color"

	qrcode "github.com/skip2/go-qrcode"
)

func main() {

	b, n := CreateQr("12324")
	if b {
		fmt.Println(n)
	}

}

func CreateQr(key string) (bool, string) {

	url := "https://www.redigo.cl/cotizacion/"
	urlqr := fmt.Sprintf("%v/%v", url, key)

	q, err := qrcode.New(urlqr, qrcode.Medium)
	if err != nil {
		return false, ""
	}

	imagename := fmt.Sprintf("%vqr.png", key)
	name := fmt.Sprintf("./tmp/%v", imagename)

	q.DisableBorder = true
	q.BackgroundColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	q.ForegroundColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}

	err = q.WriteFile(128, name)
	if err != nil {
		return false, ""
	}
	return true, name
}
