from machine import Pin
from time import sleep_ms
from machine import UART
import network

ap = network.WLAN(network.AP_IF) # create access-point interface
ap.active(False)         # activate the interface

uart = UART(0, baudrate=115200)

# while True:
uart.write("Butt Mover Connected\r\n")

button = Pin(5, Pin.OUT)

trigger = True

while True:
    if not button.value() and (not trigger):
        uart.write("Button Released\r\n")
        sleep_ms(1000)
        trigger = True

    if button.value() and trigger:
        uart.write("Button Pressed\r\n")
        sleep_ms(1000)
        trigger = False


