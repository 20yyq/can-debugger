# can-debugger
golang network can bus

### Install Candebugger
	git clone https://github.com/20yyq/can-debugger.git
	cd can-debugger
	make install

### Examples Candebugger
	./candebugber help
	./candebugber can0 read
	./candebugber vcan0 read
	./candebugber can0 write -id 12345 -data string -ext=true
	./candebugber can0 write -id 12345 -data string -ext=1
	./candebugber can0 write -id 12345 -data string -ext=True

	./candebugber vcan0 write -id 2047 -data string
	./candebugber vcan0 write -id 2047 -data string
	./candebugber vcan0 write -id 2047 -data string
	
	# fd
	./candebugber vcan0 write -id 123 -data string -flags 16 -fd=true -ext=0
	./candebugber vcan0 write -id 12345 -data string -flags 16 -fd=0 -ext=1
	
	# iface 设备开启关闭
	./candebugber vcan0 iface -bitrate 1250000 -up=true
	./candebugber vcan0 iface -up=F
	# 直接关闭设备
	./candebugber vcan0 iface