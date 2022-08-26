package releaser

import (
	"fmt"
	"nerijusdu/release-button/internal/api"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/util"
)

type Releaser struct {
	argoApi *api.ArgoApi
	configs *config.Config
}

func NewReleaser(argoApi *api.ArgoApi, configs *config.Config) *Releaser {
	return &Releaser{
		argoApi: argoApi,
		configs: configs,
	}
}

func (r *Releaser) Listen(clickChan <-chan string) {
	for button := range clickChan {
		switch button {
		case "release":
			r.Sync()
		}
	}
}

func (r *Releaser) Sync() {
	apps, err := r.argoApi.GetApps(r.configs.Selectors, false)
	if err != nil {
		fmt.Printf("ERR: Failed to get apps: Error: %v\n", err)
		return
	}

	for _, app := range apps.Items {
		if app.Status.Sync.Status == "OutOfSync" {
			fmt.Println(app.Metadata.Name + " is out of sync")
			if util.Contains(r.configs.Ignore, app.Metadata.Name) {
				fmt.Println("Skipping")
				continue
			}

			err = r.argoApi.Sync(app.Metadata.Name)
			if err != nil {
				fmt.Printf("ERR: Failed to sync %s. Error: %v\n", app.Metadata.Name, err)
			} else {
				fmt.Println("Synced " + app.Metadata.Name)
			}
		}
	}
}