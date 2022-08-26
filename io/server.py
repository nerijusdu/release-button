from gpiozero import Button
from dotenv import load_dotenv
import requests
import os
import signal
import sys

load_dotenv()
url = os.environ.get('RELEASER_URL')
dummy_io = os.environ.get('DUMMY_IO') == 'true'

def notify_releaser():
  requests.post(url+'/io/buttons/release')

if dummy_io == False:
  print('Using rpio')
  button = Button(4)
  button.when_pressed = notify_releaser


print('Listening to IO')

def sigintHandler(signal_number, frame):
  print('Exiting')
  sys.exit()

signal.signal(signal.SIGINT, sigintHandler)
signal.pause()

exit(0)