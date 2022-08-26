from gpiozero import Button
from dotenv import load_dotenv
from signal import pause
import requests
import os

load_dotenv()
url = os.environ.get('RELEASER_URL')
dummy_io = os.environ.get('DUMMY_IO') == 'true'

def notify_releaser():
  requests.post(url+'/io/button/4')

if dummy_io:
  button = Button(4)
  button.when_pressed = notify_releaser

print('Listening to IO')

pause()