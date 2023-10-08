from time import sleep
import multiprocessing
import os

dummy_io = os.environ.get('DUMMY_IO') == 'true'
pinMap = {
  "release": 4,
  "button_led": 18
}
ioMap = {}
lcd = {}

if dummy_io == False:
  from gpiozero import Button
  from gpiozero import LED
  from RPLCD.i2c import CharLCD
  lcd = CharLCD('PCF8574', 0x27)
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

def lcd_write(text):
  if dummy_io:
    return

  lcd.clear()
  lcd.write_string(text)

def lcd_clear():
  if dummy_io:
    return

  lcd.clear()