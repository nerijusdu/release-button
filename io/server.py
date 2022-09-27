from dotenv import load_dotenv
load_dotenv()

import os
import signal
import sys
import requests
import releaser_api
import releaser_io

url = os.environ.get('RELEASER_URL')

def notify_releaser():
  requests.post(url+'/io/buttons/release')

def sigintHandler(signal_number, frame):
  print('Exiting')
  sys.exit()

releaser_io.listen_to_button("release", notify_releaser)

signal.signal(signal.SIGINT, sigintHandler)
    # signal.pause()

def run():
    releaser_api.start()

run()
