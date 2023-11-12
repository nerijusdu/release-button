from time import sleep
import multiprocessing
import os
import keypad_io
import release_btn_io
import lcd_io

dummy_io = os.environ.get('DUMMY_IO') == 'true'

phonePickUpButton = None
phonePutDownButton = None

if dummy_io == False:
  from gpiozero import Button
  phonePickUpButton = Button(10)
  phonePutDownButton = Button(9)


def listen_to_button(id, listener):
  global phonePickUpButton
  global phonePutDownButton

  if id == "release":
    release_btn_io.listen_to_button("release", listener)

  if id == "phone-pickup":
    phonePickUpButton.when_pressed = listener

  if id == "phone-putdown":
    phonePutDownButton.when_pressed = listener


def led_on(id):
  release_btn_io.led_on(id)

def led_off(id):
  release_btn_io.led_off(id)

def led_blink(id, onoff):
  release_btn_io.led_blink(id, onoff)

def lcd_write(lines):
  lcd_io.lcd_write(lines)

def lcd_clear():
  lcd_io.lcd_clear()

def lcd_push(line):
  lcd_io.lcd_push(line)

keypad_thread = None

def keypad_toggle(on):
  global keypad_thread

  if dummy_io:
    return

  if on == True:
    if keypad_thread != None:
      return

    keypad_thread = multiprocessing.Process(target=keypad_io.listen_to_keypad)
    keypad_thread.start()
  else:
    if keypad_thread == None:
      return

    keypad_thread.terminate()
    keypad_thread = None