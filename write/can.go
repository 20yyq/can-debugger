// @@
// @ Author       : Eacher
// @ Date         : 2023-09-06 15:25:10
// @ LastEditTime : 2023-09-11 08:12:35
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
	flags		uint
	data		string
	ext			bool
	remote		bool
	fd			bool
)

func InitFlagArge(f *flag.FlagSet) {
	f.UintVar(&id, "id", 0, "id usage uint")
	f.UintVar(&flags, "flags", 0, "flags usage uint")
	f.StringVar(&data, "data", "", "data usage string")
	f.BoolVar(&ext, "ext", false, "ext usage bool")
	f.BoolVar(&remote, "remote", false, "remote usage bool")
	f.BoolVar(&fd, "fd", false, "fd usage bool")
	f.Usage = help
}

func Run(c *sockcan.Can) {
	var b [can.CanFDDataLength]byte
	l := len(data)
	fmt.Println("Running...")
	if l > can.CanFDDataLength {
		l = can.CanFDDataLength
	}
	if !fd && l > can.CanDataLength {
		l = can.CanDataLength
	}
	copy(b[:], data[:l])
	f := can.Frame{Len: uint8(l), CanFd: fd, Extended: ext, Remote: remote, Data: b, Flags: uint8(flags)}
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