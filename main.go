package main

import (
	"nerijusdu/release-button/pkg/api"
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
	argoApi := api.NewArgoApi()
	err := argoApi.LoadToken(api.AuthRequest{
		Username: os.Getenv("ARGOCD_USERNAME"),
		Password: os.Getenv("ARGOCD_PASSWORD"),
	})
	if err != nil {
		panic(err)
	}

	err = argoApi.Sync("test")
	if err != nil {
		panic(err)
	}
}
