from time import sleep
import multiprocessing
import os

dummy_io = os.environ.get('DUMMY_IO') == 'true'
disable_button = os.environ.get('DISABLE_RELEASE_BUTTON') == 'true'

pinMap = {
  "release": 4,
  "button_led": 18
}
ioMap = {
  "release": None,
  "button_led": None
}

if dummy_io == False and disable_button == False:
  from gpiozero import LED
  from gpiozero import Button

  button_led = LED(pinMap['button_led'])
  ioMap.button_led = button_led
  button = Button(pinMap['release'])
  ioMap.release = button


def listen_to_button(id, listener):
  if dummy_io or disable_button:
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

def _led_blink(id):
  if id not in ioMap.keys():
    return

  while True:
    ioMap[id].on()
    sleep(1)
    ioMap[id].off()
    sleep(1)

thread = None

def led_blink(id, onoff):
  if dummy_io:
    return

  global thread

  if onoff == True:
    if thread != None:
      return

    thread = multiprocessing.Process(target=_led_blink, args=(id,))
    thread.start()
  else:
    if thread == None:
      return

    thread.terminate()
    thread = None
    ioMap[id].off()
