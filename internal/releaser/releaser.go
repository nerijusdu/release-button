package releaser

import (
	"fmt"
	"nerijusdu/release-button/internal/argoApi"
	"nerijusdu/release-button/internal/config"
	"nerijusdu/release-button/internal/controls"
	"nerijusdu/release-button/internal/util"
	"time"
)

type Releaser struct {
	argoApi      argoApi.IArgoApi
	ioController *controls.IOController
	configs      *config.Config
	isSyncing    bool
}

func NewReleaser(
	aApi argoApi.IArgoApi,
	ioController *controls.IOController,
	configs *config.Config,
) *Releaser {
	return &Releaser{
		argoApi:      aApi,
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
				r.ioController.WriteToLCD([]string{
					"Failed to check status",
					err.Error(),
				})
			}
			if !inSync {
				err = r.ioController.TurnOnLed("button_led")
				if err != nil {
					fmt.Printf("ERR: turn on led. %v\n", err)
				}
			} else {
				err = r.ioController.TurnOffLed("button_led")
				if err != nil {
					fmt.Printf("ERR: turn off led. %v\n", err)
				}
			}
		},
	)

	for button := range clickChan {
		switch button {
		case "release":
			err := r.Sync(nil)
			if err != nil {
				fmt.Printf("ERR: failed to sync. %v", err)
				r.ioController.WriteToLCD([]string{
					"Oopsie woopsie",
					err.Error(),
				})
			}
		}
	}
}

func (r *Releaser) Sync(index *int) error {
	r.ioController.WriteToLCD([]string{"Vazhiojem..."})

	apps, err := r.argoApi.GetApps(r.configs.Selectors, false)
	if err != nil {
		return err
	}

	err = r.ioController.BlinkLed("button_led", true)
	if err != nil {
		return err
	}

	r.isSyncing = true

	for i, app := range apps.Items {
		if index != nil && *index != i+1 {
			continue
		}

		if app.Status.Sync.Status == "OutOfSync" {
			fmt.Println(app.Metadata.Name + " is out of sync")
			if !util.Contains(r.configs.Allowed, app.Metadata.Name) {
				fmt.Println("Skipping")
				continue
			}

			err = r.argoApi.Sync(app.Metadata.Name)
			if err != nil {
				fmt.Printf("ERR: Failed to sync %s. Error: %v\n", app.Metadata.Name, err)
				r.ioController.PushToLcd(fmt.Sprintf("Failed to sync %s", app.Metadata.Name))
			} else {
				fmt.Println("Synced " + app.Metadata.Name)
				r.ioController.PushToLcd(fmt.Sprintf("Synced %s", app.Metadata.Name))
			}
		} else if index != nil {
			fmt.Println("Already in sync")
			r.ioController.PushToLcd(fmt.Sprintf("%s is already in sync", app.Metadata.Name))
		}
	}

	util.ScheduleControlled(
		5*time.Second,
		func() bool {
			inSync, err := r.IsInSync()
			if err != nil {
				fmt.Printf("ERR: failed to check status. %v\n", err)
				r.ioController.WriteToLCD([]string{
					"Failed to check status",
					err.Error(),
				})
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
		if app.Status.Sync.Status == "OutOfSync" &&
			util.Contains(r.configs.Allowed, app.Metadata.Name) {
			return false, nil
		}
	}

	return true, nil
}
