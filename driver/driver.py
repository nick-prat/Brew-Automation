import threading
from time import sleep
from tools.button import Button
import LCD1602
import RPi.GPIO as GPIO

from tools.toggle import Toggle

FILE = "/sys/bus/w1/devices/28-3ce1d44344a1/temperature"

def c_to_f(c):
    return c * (9.0/5.0) + 32

IN1_PIN = 40
IN2_PIN = 38
COOL_PIN = 8
WARM_PIN = 10

class Controller:
    def __init__(self) -> None:
        GPIO.setmode(GPIO.BOARD)
        
        LCD1602.init(0x27, 1)
        LCD1602.write(0, 0, "0")

        self.cool_relay = Toggle(COOL_PIN)
        self.warm_relay = Toggle(WARM_PIN)

        self.up_button = Button(IN1_PIN, on_down=self.on_temp_up) 
        self.down_button = Button(IN2_PIN, on_down=self.on_temp_down)

        self.target_temp = 72
        self.previous_temp = 72

        self.read_temp = 0
        self.previous_read_temp = 0 

        self.delay = 0.01
        self.temp_thread = threading.Thread(target=self.read_temp_thread)
        self.temp_thread.start()
        self.temp_thread_event = threading.Event()
    
    def shutdown(self):
        GPIO.cleanup()
        LCD1602.clear()
        self.temp_thread_event.set()
        self.temp_thread.join()
    
    def on_temp_down(self):
        self.update_target_temp(-1)

    def on_temp_up(self):
        self.update_target_temp(1)
    
    def update_target_temp(self, diff):
        if (self.target_temp != self.previous_temp):
            return
    
        self.target_temp = self.target_temp + diff
    
    def read_temp_thread(self):
        with open(FILE) as temp_file:
            while not self.temp_thread_event.is_set():
                lines = temp_file.readlines()
                if lines:
                    self.read_temp = c_to_f(int(lines[0]) / 1000)
                temp_file.seek(0)
    
    def process(self):
        has_written = False

        if self.target_temp != self.previous_temp:
            self.previous_temp = self.target_temp
            has_written = True

        if self.read_temp != self.previous_read_temp:
            self.previous_read_temp = self.read_temp 
            has_written = True

        if has_written: 
            LCD1602.clear()
            LCD1602.write(0, 0, str(self.target_temp))
            LCD1602.write(0, 1, "{:.2f}".format(self.read_temp))
        
        if self.read_temp > self.target_temp:
            self.warm_relay.set_off()
            self.cool_relay.set_on()
        elif self.read_temp < self.target_temp:
            self.warm_relay.set_on()
            self.cool_relay.set_off()
        else:
            self.warm_relay.set_off()
            self.cool_relay.set_off()

    def loop(self):
        while True:
            self.process()
            sleep(self.delay)


if __name__ == "__main__":
    c = Controller()
    try:
        c.loop()
    except KeyboardInterrupt:
        print("Shutting down...")
        c.shutdown()