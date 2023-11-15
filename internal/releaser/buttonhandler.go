package releaser

import (
	"fmt"
	"nerijusdu/release-button/internal/controls"
)

func (r *Releaser) handleButtonClick(clickChan <-chan controls.Action) {
	cancelChan := make(chan bool)

	for a := range clickChan {
		switch a.Action {
		case "release":
			go func(action controls.Action) {
				fmt.Printf("Releasing %d \n", action.Data.Number)
				var err error
				if action.Data.Number == 0 {
					err = r.Sync(nil)
				} else {
					cancelChan = make(chan bool)
					err = r.SyncWithAudioConfirm(action.Data.Number, cancelChan)
				}

				if err != nil {
					fmt.Printf("ERR: failed to sync. %v", err)
					r.ioController.WriteToLCD([]string{
						"Oopsie woopsie",
						err.Error(),
					})
				}
			}(a)
		case "cancel":
			go func() {
				fmt.Println("Cancelling")
				cancelChan <- true
			}()
		}
	}
}
