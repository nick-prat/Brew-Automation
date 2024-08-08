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