// @@
// @ Author       : Eacher
// @ Date         : 2023-09-06 15:25:10
// @ LastEditTime : 2023-09-07 14:10:41
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  : 
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/write/can.go
// @@
package write

import (
	"fmt"
	"flag"
	
	"github.com/20yyq/can/sockcan"
	"github.com/20yyq/packet/can"
)

const helpOutput = `
	candebugber CAN-frames via CAN_RAW sockets.

	Usage: candebugber <device> <debugber_name>.

	<device>:
		Examples: 
			can0

	<debugber_name>:
		{read}:
		
		
		{write}:
		{-id}		for 'classic' CAN 2.0 data frames
			3 (SFF) or 8 (EFF) hex chars
		{-data}		for 'classic' CAN 2.0 data frames
			0-8 byte string

			Examples:
				-id 111888 -data 1122334455667788
				-id 3224123 -data DEADBEEF
				-id 11999 -data 
				-id 22123 -data 1
				-id 687213 -data 311223344

`

var (
	id			uint
	data		string
	ext			int
	remote		int
)

func InitFlagArge(f *flag.FlagSet) {
	f.UintVar(&id, "id", 0, "id usage uint")
	f.StringVar(&data, "data", "", "data usage string")
	f.IntVar(&ext, "ext", 0, "ext usage int")
	f.IntVar(&remote, "remote", 0, "remote usage int")
	f.Usage = help
}

func Run(c *sockcan.Can) {
	var b [can.DataLength]byte
	l := len(data)
	fmt.Println("Running...")
	if l > can.DataLength {
		l = can.DataLength
	}
	copy(b[:], data[:l])
	f := can.Frame{DLC: uint8(l), Extended: !(ext == 0), Remote: !(remote == 0), Data: b}
	err := f.SetID(uint32(id))
	if err != nil {
		fmt.Println("SetID err: ", err)
		return
	}
	if err = c.WriteFrame(f); err != nil {
		fmt.Println("Write err: ", err)
	}
}

func help() {
	fmt.Print(helpOutput)
}