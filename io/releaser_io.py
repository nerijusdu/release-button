from gpiozero import Button
from gpiozero import LED
import os

dummy_io = os.environ.get('DUMMY_IO') == 'true'
pinMap = {
  "release": 4,
  "button_led": 18
}
ioMap = {}

if dummy_io == False:
  button_led = LED(pinMap['button_led'])
  button = Button(pinMap['release'])
  ioMap = {
    "release": button,
    "button_led": button_led
  }


def listen_to_button(id, listener):
  if dummy_io:
    return

  if id in ioMap.keys():
    ioMap[id].when_pressed = listener


def led_on(id):
  if dummy_io:
    return

  if id in ioMap.keys():
    ioMap[id].on()


def led_off(id):
  if dummy_io:
    return

  if id in ioMap.keys():
    ioMap[id].off()