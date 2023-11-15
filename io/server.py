from dotenv import load_dotenv
load_dotenv()

import os
import signal
import sys
import requests
import releaser_api
import releaser_io
import keypad_io

url = os.environ.get('RELEASER_URL')

def notify_releaser():
  requests.post(url+'/io/actions/release')

def notify_keypad(num):
  requests.post(url+'/io/actions/release', json={'number': int(num)})

def notify_cancel():
  requests.post(url+'/io/actions/cancel')

def putdown():
  releaser_io.keypad_toggle(False)
  notify_cancel()

def pickup():
  releaser_io.keypad_toggle(True)

def sigintHandler(signal_number, frame):
  print('Exiting')
  sys.exit()

releaser_io.listen_to_button("release", notify_releaser)
releaser_io.listen_to_button("phone-pickup", pickup)
releaser_io.listen_to_button("phone-putdown", putdown)
keypad_io.setup(notify_keypad, notify_cancel)

signal.signal(signal.SIGINT, sigintHandler)
    # signal.pause()

def run():
  releaser_api.start()

run()
