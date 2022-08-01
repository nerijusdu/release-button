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
		apps, err := argoApi.GetApps(c.Selectors, false)
		if err != nil {
			panic(err)
		}

		for _, app := range apps.Items {
			if app.Status.Sync.Status == "OutOfSync" {
				fmt.Println(app.Metadata.Name + " is out of sync")
				// argoApi.Sync(app.Metadata.Name)
				for k, v := range app.Metadata.Labels {
					fmt.Println(k + ": " + v)
				}

				fmt.Println("")
			}
		}
		fmt.Println("Done processing click")
	}
}
