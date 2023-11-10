package main

import (
	"fmt"
	"nerijusdu/release-button/internal/argoApi"
	"nerijusdu/release-button/internal/audio"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/controls"
	"nerijusdu/release-button/internal/releaser"
	"nerijusdu/release-button/internal/web"

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
	clickChan := make(chan string)
	ioListener := controls.NewIOListener()
	ioController := controls.NewIOController()
	synth := audio.NewSynth()
	releaser := releaser.NewReleaser(aApi, ioController, synth, c)
	webApi := web.NewWebApi(aApi, c)

	go webApi.Listen()
	go ioListener.Listen(clickChan)
	releaser.Listen(clickChan)
}
