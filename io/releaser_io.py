from time import sleep
import multiprocessing
import os
import keypad_io

dummy_io = os.environ.get('DUMMY_IO') == 'true'
disable_lcd = os.environ.get('DISABLE_LCD') == 'true'
disable_button = os.environ.get('DISABLE_BUTTON') == 'true'
disable_button_led = os.environ.get('DISABLE_BUTTON_LED') == 'true'

pinMap = {
  "release": 4,
  "button_led": 18
}
ioMap = {
  "release": None,
  "button_led": None
}
lcd = {}

if dummy_io == False:
  if disable_lcd == False:
    from RPLCD.i2c import CharLCD
    lcd = CharLCD('PCF8574', 0x27, auto_linebreaks=False)
  if disable_button_led == False:
    from gpiozero import LED
    button_led = LED(pinMap['button_led'])
    ioMap.button_led = button_led
  if disable_button == False:
    from gpiozero import Button
    button = Button(pinMap['release'])
    ioMap.release = button


def listen_to_button(id, listener):
  if dummy_io or disable_button:
    return

  if id in ioMap.keys():
    ioMap[id].when_pressed = listener


def led_on(id):
  if dummy_io:
    return

  if id in ioMap.keys():
    ioMap[id].on()


def led_off(id):
  if dummy_io:
    return

  if id in ioMap.keys():
    ioMap[id].off()

def _led_blink(id):
  if id not in ioMap.keys():
    return

  while True:
    ioMap[id].on()
    sleep(1)
    ioMap[id].off()
    sleep(1)

thread = None

def led_blink(id, onoff):
  global thread

  if onoff == True:
    if thread != None:
      return

    thread = multiprocessing.Process(target=_led_blink, args=(id,))
    thread.start()
  else:
    if thread == None:
      return

    thread.terminate()
    thread = None
    ioMap[id].off()

lcd_thread = None

def _lcd_write(line, line_num):
  s = line + '   '

  while True:
    for i in range(len(s) - 16 + 1):
      lcd.cursor_pos = (line_num, 0)
      lcd.write_string(s[i:i+16])
      sleep(0.5)

def lcd_write(lines):
  global lcd_thread

  if dummy_io:
    return

  lcd_clear()

  for i in range(len(lines)):
    if len(lines[i]) > 16:
      lcd_thread = multiprocessing.Process(target=_lcd_write, args=(lines[i],i))
      lcd_thread.start()
    else:
      lcd.cursor_pos = (i, 0)
      lcd.write_string(lines[i])

def lcd_clear():
  global lcd_thread

  if dummy_io:
    return

  if lcd_thread != None:
    lcd_thread.terminate()
    lcd_thread = None

  lcd.clear()

lcd_buffer = []
lcd_buffer_thread = None

def _lcd_push():
  global lcd_buffer
  global lcd_buffer_thread

  while len(lcd_buffer) > 0:
    lcd_write(lcd_buffer)
    lcd_buffer.pop(0)
    sleep(1)

  lcd_buffer_thread.terminate()
  lcd_buffer_thread = None

def lcd_push(line):
  if dummy_io:
    return

  lcd_buffer.append(line)

  if lcd_buffer_thread == None:
    lcd_buffer_thread = multiprocessing.Process(target=_lcd_push)
    lcd_buffer_thread.start()

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