package controls

import (
	"bytes"
	"encoding/json"
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

type writeLcdRequest struct {
	Text []string `json:"text"`
}

type pushLcdRequest struct {
	Text string `json:"text"`
}

func (c *IOController) WriteToLCD(lines []string) error {
	data, err := json.Marshal(&writeLcdRequest{Text: lines})
	if err != nil {
		return err
	}

	_, err = http.Post(c.url+"/io/lcd/write", "application/json", bytes.NewReader(data))
	return err
}

func (c *IOController) PushToLcd(line string) error {
	data, err := json.Marshal(&pushLcdRequest{Text: line})
	if err != nil {
		return err
	}

	_, err = http.Post(c.url+"/io/lcd/push", "application/json", bytes.NewReader(data))
	return err
}

type speakRequest struct {
	Text string `json:"text"`
}

func (c *IOController) Speak(text string) error {
	data, err := json.Marshal(&speakRequest{Text: text})
	if err != nil {
		return err
	}

	_, err = http.Post(c.url+"/io/synth", "application/json", bytes.NewReader(data))
	return err
}
