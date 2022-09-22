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
	argoApi      *api.ArgoApi
	ioController *controls.IOController
	configs      *config.Config
	isSyncing    bool
}

func NewReleaser(
	argoApi *api.ArgoApi,
	ioController *controls.IOController,
	configs *config.Config,
) *Releaser {
	return &Releaser{
		argoApi:      argoApi,
		ioController: ioController,
		configs:      configs,
	}
}

func (r *Releaser) Listen(clickChan <-chan string) {
	util.Schedule(
		time.Duration(r.configs.RefreshInterval)*time.Second,
		func() {
			if r.isSyncing {
				return
			}

			inSync, err := r.IsInSync()
			if err != nil {
				fmt.Printf("ERR: failed to check status. %v\n", err)
			}
			if !inSync {
				err = r.ioController.TurnOnLed("button_led")
				if err != nil {
					fmt.Printf("ERR: turn on led. %v\n", err)
				}
			}
		},
	)

	for button := range clickChan {
		switch button {
		case "release":
			err := r.Sync()
			if err != nil {
				fmt.Printf("ERR: failed to sync. %v", err)
			}
		}
	}
}

func (r *Releaser) Sync() error {

	apps, err := r.argoApi.GetApps(r.configs.Selectors, false)
	if err != nil {
		return err
	}

	err = r.ioController.BlinkLed("button_led", true)
	if err != nil {
		return err
	}

	r.isSyncing = true

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

	util.ScheduleControlled(
		5*time.Second,
		func() bool {
			inSync, err := r.IsInSync()
			if err != nil {
				fmt.Printf("ERR: failed to check status. %v\n", err)
			}
			if inSync {
				err = r.ioController.BlinkLed("button_led", false)
				if err != nil {
					fmt.Printf("ERR: failed to turn off led. %v", err)
				}
				r.isSyncing = false
				return true
			}
			return false
		},
	)

	return nil
}

func (r *Releaser) IsInSync() (bool, error) {
	apps, err := r.argoApi.GetApps(r.configs.Selectors, true)
	if err != nil {
		return true, err
	}

	for _, app := range apps.Items {
		if app.Status.Sync.Status == "OutOfSync" && !util.Contains(r.configs.Ignore, app.Metadata.Name) {
			return false, nil
		}
	}

	return true, nil
}
