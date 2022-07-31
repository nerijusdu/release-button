package main

import (
	"fmt"
	"nerijusdu/release-button/internal/api"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/controls"
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

	clickChan := make(chan bool)
	cont := controls.NewController()
	go cont.WaitForClick("GPIO4", clickChan)
	fmt.Println("Waiting for clicks")
	for range clickChan {
		fmt.Println("Got click")
		for _, v := range c.Apps {
			err = argoApi.Sync(v)
			if err != nil {
				panic(err)
			}
		}
	}
}
