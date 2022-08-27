import io
from flask import Flask
import releaser_io

app = Flask(__name__)

@app.route("/health")
def health():
  return "OK"


@app.route("/io/led/<id>/on", methods=['POST'])
def led_on(id):
  releaser_io.led_on(id)
  return ""

@app.route("/io/led/<id>/off", methods=['POST'])
def led_on(id):
  releaser_io.led_off(id)
  return ""

def start():
  app.run(host='0.0.0.0', port=6968)