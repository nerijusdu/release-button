from gpiozero import Button
from gpiozero import LED
import os

dummy_io = os.environ.get('DUMMY_IO') == 'true'
pinMap = {
  "release": 4,
  "button_led": 12
}

def listen_to_button(pin, listener):
  if dummy_io:
    return

  button = Button(pinMap[pin])
  button.when_pressed = listener


def led_on(pin):
  if dummy_io:
    return

  led = LED(pinMap[pin])
  led.on()


def led_off(pin):
  if dummy_io:
    return

  led = LED(pinMap[pin])
  led.off()