package main

import (
	"fmt"
	"nerijusdu/release-button/internal/argoApi"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/controls"
	"nerijusdu/release-button/internal/releaser"
	"nerijusdu/release-button/internal/web"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Failed to load .env file. %v \n", err)
	}
}

func main() {
	fmt.Println("Starting")
	defer fmt.Println("Exiting")

	c, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	aApi := argoApi.NewArgoApi()
	err = aApi.LoadToken(argoApi.AuthRequest{
		Username: os.Getenv("ARGOCD_USERNAME"),
		Password: os.Getenv("ARGOCD_PASSWORD"),
	})
	if err != nil {
		fmt.Printf("Failed to fetch argo token. %v \n", err)
	}

	clickChan := make(chan string)
	ioListener := controls.NewIOListener()
	ioController := controls.NewIOController()
	releaser := releaser.NewReleaser(aApi, ioController, c)
	webApi := web.NewWebApi(aApi, c)

	go webApi.Listen()
	go ioListener.Listen(clickChan)
	releaser.Listen(clickChan)
}
