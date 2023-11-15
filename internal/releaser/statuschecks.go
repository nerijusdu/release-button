package releaser

import (
	"fmt"
	"nerijusdu/release-button/internal/util"
)

func (r *Releaser) checkIfNeedsSyncing() {
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
}

func (r *Releaser) updateSyncProgress() bool {
	inSync, err := r.IsInSync()
	if err != nil {
		fmt.Printf("ERR: failed to check status. %v\n", err)
		r.ioController.WriteToLCD([]string{
			"Failed to check status",
			err.Error(),
		})
		return false
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
