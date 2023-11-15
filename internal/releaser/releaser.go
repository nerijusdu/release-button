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

func (r *Releaser) Listen(clickChan <-chan controls.Action) {
	util.Schedule(
		time.Duration(r.configs.RefreshInterval)*time.Second,
		r.checkIfNeedsSyncing,
	)

	r.handleButtonClick(clickChan)
}

func (r *Releaser) Sync(index *int) error {
	r.ioController.WriteToLCD([]string{"Vazhiojem..."})
	r.ioController.Speak("Vahzioyam")

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

	// util.ScheduleControlled(
	// 	5*time.Second,
	// 	r.updateSyncProgress,
	// )

	return nil
}

func (r *Releaser) SyncWithAudioConfirm(index int, cancelChan <-chan bool) error {
	apps, err := r.argoApi.GetApps(r.configs.Selectors, false)
	if err != nil {
		r.ioController.Speak("Failed to get services")
		return err
	}

	var length int = len(apps.Items)
	if index > length {
		r.ioController.Speak("Invalid index")
		return err
	}

	app := apps.Items[index-1]

	r.ioController.Speak(fmt.Sprintf("Releasing %s in 5 seconds", app.Metadata.Name))

	select {
	case <-time.After(5 * time.Second):
		err = r.Sync(&index)
		if err != nil {
			r.ioController.Speak("Failed to release")
			return err
		}
	case <-cancelChan:
		r.ioController.Speak("Cancelled")
		fmt.Println("Sync with audio cancelled")
		return nil
	}

	return nil
}
