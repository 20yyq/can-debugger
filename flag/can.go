// @@
// @ Author       : Eacher
// @ Date         : 2023-09-06 13:56:59
// @ LastEditTime : 2023-09-07 14:11:10
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  : 
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/flag/can.go
// @@
package flag

import (
	"fmt"
	"os"
	"flag"
)

type FlagSetFunc func(*flag.FlagSet)

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
				-id 888 -data 1122334455667788
				-id 123 -data DEADBEEF
				-id 999 -data 
				-id 123 -data 1
				-id 213 -data 311223344

`

var (
	canInterfaceName string
	canDebuggerName string
)

func init() {
	if len(os.Args) < 3 {
		if len(os.Args) > 0 && (os.Args[1] == "help" || os.Args[1] == "-h") {
			fmt.Print(helpOutput)
			os.Exit(0)
		}
		fmt.Println("os.Args need min three arg")
		os.Exit(1)
	}
	canInterfaceName, canDebuggerName = os.Args[1], os.Args[2]
}

func Init(f FlagSetFunc) error {
	fs := flag.NewFlagSet(canDebuggerName, flag.ContinueOnError)
	f(fs)
	return fs.Parse(os.Args[3:])
}

func InterfaceName() string {
	return canInterfaceName
}

func DebuggerName() string {
	return canDebuggerName
}