package releaser

import (
	"fmt"
	"nerijusdu/release-button/internal/controls"
)

func (r *Releaser) handleButtonClick(clickChan <-chan controls.Action) {
	cancelChan := make(chan bool, 1)

	for a := range clickChan {
		switch a.Action {
		case "release":
			var err error
			if a.Data.Number == 0 {
				err = r.Sync(nil)
			} else {
				err = r.SyncWithAudioConfirm(a.Data.Number, cancelChan)
			}

			if err != nil {
				fmt.Printf("ERR: failed to sync. %v", err)
				r.ioController.WriteToLCD([]string{
					"Oopsie woopsie",
					err.Error(),
				})
			}
		case "cancel":
			fmt.Println("Cancelling")
			cancelChan <- true
			fmt.Println("Cancelled")
		}
	}
}
