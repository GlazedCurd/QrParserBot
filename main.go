package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	tele "gopkg.in/telebot.v3"
)

func buildBot() (*tele.Bot, error) {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("env BOT_TOKEN is empty; No token found")
	}
	webHookHost := os.Getenv("BOT_PUBLIC_HOST")
	if webHookHost == "" {
		return nil, fmt.Errorf("env BOT_PUBLIC_HOST is empty; No webhook host found")
	}
	webHookPath := os.Getenv("BOT_PATH")
	if webHookPath == "" {
		return nil, fmt.Errorf("env BOT_PATH is empty; No webhook path found")
	}
	webHookPort := os.Getenv("BOT_PUBLIC_PORT")
	if webHookPort == "" {
		webHookPort = "443"
	}

	publicCertPath := os.Getenv("BOT_PUBLIC_CERT_PATH")
	privateCertPath := os.Getenv("BOT_PRIVATE_CERT_PATH")

	debug := os.Getenv("DEBUG")
	verbose := debug != ""

	wh := &tele.Webhook{
		Endpoint: &tele.WebhookEndpoint{
			PublicURL: fmt.Sprintf("%s:%s%s", webHookHost, webHookPort, webHookPath),
		},
		AllowedUpdates: []string{"message"},
		Listen:         ":" + webHookPort,
	}
	if publicCertPath != "" && privateCertPath != "" {
		wh.Endpoint.Cert = publicCertPath
		wh.TLS = &tele.WebhookTLS{
			Key:  privateCertPath,
			Cert: publicCertPath,
		}
	}

	settings := tele.Settings{
		Token:   token,
		Verbose: verbose,
		Poller:  wh,
	}

	return tele.NewBot(settings)
}

func main() {
	b, err := buildBot()
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tele.OnPhoto, func(c tele.Context) error {
		if c.Message().Photo == nil {
			fmt.Println("Photo is nil")
			return nil
		}
		photoReader, err := b.File(&c.Message().Photo.File)
		if err != nil {
			fmt.Printf("Error on getting image %s \n", err)
		}
		img, _, err := image.Decode(photoReader)
		if err != nil {
			fmt.Printf("Error on decoding %s \n", err)
			return nil
		}
		bmp, err := gozxing.NewBinaryBitmapFromImage(img)
		if err != nil {
			log.Fatalf("Failed to make bitmap %s", err)
		}

		// decode image
		qrReader := qrcode.NewQRCodeReader()
		result, err := qrReader.Decode(bmp, nil)

		if err != nil {
			log.Fatalf("Failed to decode qr %s", err)
		}
		return c.Send(result.GetText())
	})

	b.Start()
}
