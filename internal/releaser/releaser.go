package releaser

import (
	"fmt"
	"nerijusdu/release-button/internal/api"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/controls"
	"nerijusdu/release-button/internal/util"
	"time"
)

type Releaser struct {
	argoApi *api.ArgoApi
	ioController *controls.IOController
	configs *config.Config
}

func NewReleaser(
	argoApi *api.ArgoApi,
	ioController *controls.IOController,
	configs *config.Config,
) *Releaser {
	return &Releaser{
		argoApi: argoApi,
		ioController: ioController,
		configs: configs,
	}
}

func (r *Releaser) Listen(clickChan <-chan string) {
	util.Schedule(5*time.Second, func() { 
		err := r.ChcekStatus() 
		fmt.Printf("ERR: failed to check status. %v\n", err)
	})

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

func (r *Releaser) ChcekStatus() error {
	apps, err := r.argoApi.GetApps(r.configs.Selectors, true)
	if err != nil {
		return err
	}

	for _, app := range apps.Items {
		if app.Status.Sync.Status == "OutOfSync" {
			return r.ioController.TurnOnLed("button_led")
		}
	}

	return r.ioController.TurnOffLed("button_led")
}