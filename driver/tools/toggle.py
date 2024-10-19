import RPi.GPIO as GPIO

class Toggle:
    def __init__(self, pin):
        self.pin = pin
        self.is_on = 0
        GPIO.setup(self.pin, GPIO.OUT)
        GPIO.output(self.pin, 0)
    
    def toggle(self):
        GPIO.output(self.pin, 0 if self.is_on else 1)
        self.is_on = not self.is_on 

    def set_on(self):
        if not self.is_on:
            self.is_on = True
            GPIO.output(self.pin, 1)
    
    def set_off(self):
        if self.is_on:
            self.is_on = False
            GPIO.output(self.pin, 0)
    
class ThreeWaySwitch:
    def __init__(self, high_pin, low_pin):
        self.high_pin = high_pin
        self.low_pin = low_pin
        self.state = 0

        GPIO.setup(self.low_pin, GPIO.OUT)
        GPIO.output(self.low_ping, 0)

        GPIO.setup(self.high_pin, GPIO.OUT)
        GPIO.output(self.high_pin, 0)
    
    def toggle(self):
        GPIO.output(self.pin, 0 if self.is_on else 1)
        self.is_on = not self.is_on 
    
    def set_low(self):
        GPIO.output(self.low_pin, 1)
        GPIO.output(self.high_pin, 0)
        self.state = -1

    def set_high(self):
        GPIO.output(self.low_pin, 0)
        GPIO.output(self.high_pin, 1)
        self.state = 1
    
    def set_off(self):
        GPIO.output(self.low_pin, 0)
        GPIO.output(self.high_pin, 0)
        self.state = 0
    
    def set_state(self, state):
        if self.state == 0:
            self.set_off()
        elif self.state < 0:
            self.set_low()
        else:
            self.set_high()