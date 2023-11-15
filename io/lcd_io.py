from time import sleep
import multiprocessing
import os

dummy_io = os.environ.get('DUMMY_IO') == 'true'
disable_lcd = os.environ.get('DISABLE_LCD') == 'true'
lcd = {}

if dummy_io == False:
  if disable_lcd == False:
    from RPLCD.i2c import CharLCD
    lcd = CharLCD('PCF8574', 0x27, auto_linebreaks=False)


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

  if dummy_io or disable_lcd:
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

  if dummy_io or disable_lcd:
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
  if dummy_io or disable_lcd:
    return

  lcd_buffer.append(line)

  if lcd_buffer_thread == None:
    lcd_buffer_thread = multiprocessing.Process(target=_lcd_push)
    lcd_buffer_thread.start()
