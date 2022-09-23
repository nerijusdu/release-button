from dotenv import load_dotenv
import requests
import os
import signal
import sys
import releaser_argoApi
import releaser_io

load_dotenv()
url = os.environ.get('RELEASER_URL')

def notify_releaser():
  requests.post(url+'/io/buttons/release')

releaser_io.listen_to_button("release", notify_releaser)

def sigintHandler(signal_number, frame):
  print('Exiting')
  sys.exit()

signal.signal(signal.SIGINT, sigintHandler)
# signal.pause()

releaser_argoApi.start()