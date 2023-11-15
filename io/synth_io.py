import multiprocessing
import boto3
from ctypes import *
from contextlib import contextmanager
from contextlib import closing
import os
import pyaudio

# Error handler for ALSA
ERROR_HANDLER_FUNC = CFUNCTYPE(None, c_char_p, c_int, c_char_p, c_int, c_char_p)

def py_error_handler(filename, line, function, err, fmt):
  pass

c_error_handler = ERROR_HANDLER_FUNC(py_error_handler)

@contextmanager
def noalsaerr():
  asound = cdll.LoadLibrary('libasound.so.2')
  asound.snd_lib_error_set_handler(c_error_handler)
  yield
  asound.snd_lib_error_set_handler(None)
# -----------------------

polly_client = boto3.Session(
  aws_access_key_id=os.environ.get('AWS_ACCESS_KEY_ID'),                     
  aws_secret_access_key=os.environ.get('AWS_SECRET_ACCESS_KEY'),
  region_name=os.environ.get('AWS_REGION')
).client('polly')

with noalsaerr():
  p = pyaudio.PyAudio()

# multiple requests throws an error, device is busy

synth_thread = None

def _synthesize(text):
  global synth_thread
  print('synthesizing speech')
  stream = p.open(format=p.get_format_from_width(2),
    channels=1,
    rate=16000,
    output=True)

  response = polly_client.synthesize_speech(VoiceId='Stephen',
    OutputFormat='pcm', 
    Text = text,
    Engine = 'neural',
    SampleRate = '16000')

  with closing(response["AudioStream"]) as polly_stream:
    while True:
      data = polly_stream.read(4096)
      if data is None or len(data) == 0:
        break
      stream.write(data)

  stream.stop_stream()
  stream.close()
  print('done synthesizing speech')

def synthesize_speech(text):
  global synth_thread

  if synth_thread != None:
    synth_thread.terminate()
    synth_thread = None

  synth_thread = multiprocessing.Process(target=_synthesize, args=(text,))
  synth_thread.start()


def cleanup():
  p.terminate()