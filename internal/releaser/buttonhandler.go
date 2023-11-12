package releaser

import (
	"fmt"
	"nerijusdu/release-button/internal/util"
	"strconv"
	"time"
)

func (r *Releaser) handleButtonClick(clickChan <-chan string) {
	cancelChan := make(chan bool, 1)
	cancelTimer := func() {
		if r.numInputTimerCloser != nil {
			r.numInputTimerCloser()
			r.numInputTimerCloser = nil
		}
	}

	resetTimer := func() {
		cancelTimer()
		r.numInputTimerCloser = util.Schedule(3*time.Second, func() {
			parsedNum, err := strconv.Atoi(r.numInput)
			if err != nil {
				fmt.Printf("ERR: failed to parse input: %s. %v", r.numInput, err)
				return
			}
			r.SyncWithAudioConfirm(parsedNum, cancelChan)
		})
	}

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
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			fmt.Printf("Syncing %s\n", button)
			r.numInput += button
			resetTimer()
		case "Redial", "Cancel":
			r.numInput = ""
			fmt.Println("Cancelling")
			cancelChan <- true
			fmt.Println("Cancelled")
			cancelChan = make(chan bool, 1)
			cancelTimer()
		case "R":
			r.numInput = r.numInput[:len(r.numInput)-1]
			resetTimer()
		}
	}
}
