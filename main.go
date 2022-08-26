package main

import (
	"fmt"
	"nerijusdu/release-button/internal/api"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/controls"
	"nerijusdu/release-button/internal/releaser"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Starting")
	defer fmt.Println("Exiting")

	c, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	argoApi := api.NewArgoApi()
	err = argoApi.LoadToken(api.AuthRequest{
		Username: os.Getenv("ARGOCD_USERNAME"),
		Password: os.Getenv("ARGOCD_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	clickChan := make(chan string)
	ioServer := controls.NewIOServer()
	releaser := releaser.NewReleaser(argoApi, c)

	go ioServer.Listen(clickChan)
	releaser.Listen(clickChan)
}
