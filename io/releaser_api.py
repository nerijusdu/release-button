import io
from flask import Flask
from flask import request
import releaser_io
import synth_io

app = Flask(__name__)

@app.route("/health")
def health():
  return "OK"


@app.route("/io/led/<id>/on", methods=['POST'])
def led_on_route(id):
  releaser_io.led_on(id)
  return ""

@app.route("/io/led/<id>/off", methods=['POST'])
def led_off_route(id):
  releaser_io.led_off(id)
  return ""

@app.route("/io/led/<id>/blink/<onoff>", methods=['POST'])
def led_blink_route(id, onoff):
  releaser_io.led_blink(id, onoff.lower() == "true")
  return ""

@app.route("/io/lcd/write", methods=['POST'])
def lcd_write():
  body = request.get_json()
  releaser_io.lcd_write(body['text'])
  return ""

@app.route("/io/lcd/push", methods=['POST'])
def lcd_push():
  body = request.get_json()
  releaser_io.lcd_push(body['text'])
  return ""

@app.route("/io/lcd/clear", methods=['POST'])
def lcd_clear():
  releaser_io.lcd_clear()
  return ""

@app.route("/io/keypad", methods=['POST'])
def keypad_toggle():
  body = request.get_json()
  releaser_io.keypad_toggle(body['on'])
  return ""

@app.route("/io/synth", methods=['POST'])
def synth():
  body = request.get_json()
  synth_io.synthesize_speech(body['text'])
  return ""

def start():
  print("Running IO controller on port :6968")
  app.run(host='0.0.0.0', port=6968)
