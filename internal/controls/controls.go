package controls

import (
	"fmt"
	"log"
	"os"
	"time"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

type Controller struct {
	isDummy bool
}

func NewController() *Controller {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	return &Controller{
		isDummy: os.Getenv("DUMMY_IO") == "true",
	}
}

func (c *Controller) WaitForClick(pin string, clickChan chan<- bool) {
	if c.isDummy {
		c.waitForClickDummy(clickChan)
		return
	}

	p := gpioreg.ByName("GPIO4")
	if p == nil {
		fmt.Println("Failed to find GPIO4")
		return
	}

	fmt.Printf("%s: %s\n", p, p.Function())

	if err := p.In(gpio.PullDown, gpio.BothEdges); err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		p.WaitForEdge(-1)
		clickChan <- true
	}
}

func (c *Controller) waitForClickDummy(clickChan chan<- bool) {
	for {
		time.Sleep(time.Second * 10)
		clickChan <- true
	}
}
