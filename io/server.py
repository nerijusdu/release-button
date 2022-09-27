from dotenv import load_dotenv
load_dotenv()

import os
import signal
import sys
import releaser_api
import releaser_io

def notify_releaser(url):
  requests.post(url+'/io/buttons/release')

def sigintHandler(signal_number, frame):
  print('Exiting')
  sys.exit()

url = os.environ.get('RELEASER_URL')

releaser_io.listen_to_button("release", notify_releaser)

signal.signal(signal.SIGINT, sigintHandler)
    # signal.pause()

def run():
    releaser_api.start()

run()
