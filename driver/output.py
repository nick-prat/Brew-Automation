import tools.LCD1602 as LCD1602

class ConsoleOutputStrategy:
    def process(self, read_temp, target_temp):
        print("{:.2f} - {}".format(read_temp, target_temp))


class LCDOutputStrategy:
    def process(self, read_temp, target_temp):
        LCD1602.clear()
        LCD1602.write(0, 0, str(target_temp))
        LCD1602.write(0, 1, "{:.2f}".format(read_temp))


class ThrottleOutputDecorator:
    def __init__(self, strategy):
        self.strategy = strategy
        self.prev_read_temp = None
        self.prev_target_temp = None
    
    def process(self, read_temp, target_temp):
        if read_temp != self.prev_read_temp or target_temp != self.prev_target_temp:
            self.strategy.process(read_temp, target_temp)
            self.prev_read_temp = read_temp
            self.prev_target_temp = target_temp


class MultiOutputDecorator:
    def __init__(self, *args):
        self.strategies = args
    
    def process(self, read_temp, target_temp):
        for strategy in self.strategies:
            strategy.process(read_temp, target_temp)

