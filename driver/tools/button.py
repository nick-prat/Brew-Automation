import RPi.GPIO as GPIO

class Button:
    def __init__(self, channel, on_down=None, on_up=None, bouncetime=10):
        GPIO.setup(channel, GPIO.IN, pull_up_down=GPIO.PUD_UP)
        GPIO.add_event_detect(channel, GPIO.BOTH, callback=self.process, bouncetime=bouncetime)
        self.is_pressed = False
        self.on_down = on_down
        self.on_up = on_up
    
    def process(self, channel):
        state = GPIO.input(channel)

        if not state:
            if not self.is_pressed and self.on_down:
                self.on_down()
    
            self.is_pressed = True
        
        if state:
            if self.is_pressed and self.on_up:
                self.on_up()
            
            self.is_pressed = False
    