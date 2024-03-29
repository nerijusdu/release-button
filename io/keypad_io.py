import RPi.GPIO as GPIO
from time import sleep
from threading import Timer

GPIO.setmode(GPIO.BCM)

# pins1 = [40,38,36,32,26,24] # Board
pins1 = [21,20,16,12,7,8] #BCM
# pins2 = [23,29,31,33,35,37] # Board
pins2 = [11,5,6,13,19,26] # BCM

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

onNumberEntered = None
onCancel = None

numberInput = ''
numberInputTimer = None

def setup(numberEnteredHandler, cancelHandler):
  global onNumberEntered
  global onCancel
  onNumberEntered = numberEnteredHandler
  onCancel = cancelHandler

def clearTimer():
  global numberInputTimer
  if numberInputTimer is not None:
    numberInputTimer.cancel()
    numberInputTimer = None

def sendNumber():
  global numberInput
  global onNumberEntered
  if numberInput != '':
    onNumberEntered(numberInput)
    numberInput = ''
    clearTimer()

def resetTimer():
  global numberInputTimer
  clearTimer()
  numberInputTimer = Timer(3, sendNumber)
  numberInputTimer.start()


def on_clicked(key):
  global numberInput
  global numberInputTimer
  global onNumberEntered
  global onCancel

  print('current number: ', numberInput, ' key pressed: ', key)

  if key == 'Redial':
    numberInput = ''
    clearTimer()
  elif key == 'R':
    numberInput = ''
    resetTimer()
    onCancel()
  elif key == '#':
    print('dont know what to do with #')
  elif key == '*':
    print('dont know what to do with *')
  else:
    numberInput += key
    resetTimer()

def listen_to_keypad():
  global onNumberEntered

  while True:
    for key in buttonKeys:
      GPIO.setup(buttons[key][0], GPIO.OUT)
      GPIO.setup(buttons[key][1], GPIO.IN, pull_up_down=GPIO.PUD_DOWN)
      GPIO.output(buttons[key][0], GPIO.HIGH)
      if GPIO.input(buttons[key][1]) == GPIO.HIGH:
        on_clicked(str(key))
        sleep(0.7) # avoid long press triggering multiple times
      GPIO.cleanup(buttons[key])
    sleep(0.1)