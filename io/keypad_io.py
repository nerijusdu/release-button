import RPi.GPIO as GPIO
from time import sleep

GPIO.setmode(GPIO.BOARD)

pins1 = [40,38,36,32,26,24]
pins2 = [23,29,31,33,35,37]

buttons = {
  1: [pins1[2], pins2[4]],
  4: [pins1[3], pins2[3]],
  9: [pins1[4], pins2[4]],
  2: [pins1[0], pins1[2]],
  5: [pins1[3], pins2[4]],
  3: [pins1[1], pins1[2]],
  8: [pins1[4], pins2[3]],
  6: [pins1[3], pins2[5]],
  7: [pins1[1], pins1[3]],
  0: [pins1[4], pins2[2]],
  '#': [pins1[5], pins2[1]],
  '*': [pins1[1], pins1[4]],
  'R': [pins1[5], pins2[4]],
  'Redial': [pins2[0], pins2[1]],
}

buttonKeys = list(buttons.keys())

onClickHandler = None

def setup(onClick):
  global onClickHandler
  onClickHandler = onClick

def listen_to_keypad():
  global onClickHandler

  while True:
    for key in buttonKeys:
      GPIO.setup(buttons[key][0], GPIO.OUT)
      GPIO.setup(buttons[key][1], GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
      GPIO.output(buttons[key][0], GPIO.HIGH)
      if GPIO.input(buttons[key][1]) == GPIO.HIGH:
        onClickHandler(str(key))
        sleep(0.5)
      GPIO.cleanup(buttons[key])
    sleep(0.2)