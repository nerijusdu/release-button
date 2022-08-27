package controls

import (
	"net/http"
	"strconv"
)

type IOController struct {
	url string
}

func NewIOController() *IOController {
	return &IOController{
		url: "http://localhost:6968",
	}
}

func (c *IOController) TurnOnLed(id string) error {
	_, err := http.Post(c.url+"/io/led/"+id+"/on", "application/json", nil)
	return err
}

func (c *IOController) TurnOffLed(id string) error {
	_, err := http.Post(c.url+"/io/led/"+id+"/off", "application/json", nil)
	return err
}

func (c *IOController) BlinkLed(id string, onOff bool) error {
	_, err := http.Post(
		c.url+"/io/led/"+id+"/blink/"+strconv.FormatBool(onOff), 
		"application/json", 
		nil,
	)
	return err
}