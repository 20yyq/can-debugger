// @@
// @ Author       : Eacher
// @ Date         : 2023-09-06 13:56:59
// @ LastEditTime : 2024-01-09 16:29:37
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  :
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/flag/can.go
// @@
package flag

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/20yyq/can-debugger/iface"
	"github.com/20yyq/can-debugger/read"
	"github.com/20yyq/can-debugger/sockcan"
	"github.com/20yyq/can-debugger/write"
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
	canDebuggerName  string
)

var (
	runMap map[string]func() error
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
	runMap = map[string]func() error{
		"read":  readRuning,
		"write": writeRuning,
		"iface": ifaceRuning,
	}
}

func Init() error {
	fs := flag.NewFlagSet(canDebuggerName, flag.ContinueOnError)
	switch canDebuggerName {
	case "write":
		write.InitFlagArge(fs)
	case "iface":
		iface.InitFlagArge(fs)
	case "read":
		fallthrough
	case "":
		return nil
	default:
		return fmt.Errorf("flag init '%s' not", canDebuggerName)
	}
	return fs.Parse(os.Args[3:])
}

func Runing() error {
	if f, ok := runMap[canDebuggerName]; ok && f != nil {
		return f()
	}
	return fmt.Errorf("running '%s' not method", canDebuggerName)
}

func ifaceRuning() error {
	notify, stop := make(chan struct{}), make(chan struct{})
	go listening(notify, stop)
	err := iface.Run(canInterfaceName)
	close(stop)
	<-notify
	return err
}

func readRuning() error {
	return sockcanRuning(read.Run)
}

func writeRuning() error {
	return sockcanRuning(write.Run)
}

func sockcanRuning(f func(*sockcan.Can)) error {
	can, err := sockcan.NewCan(canInterfaceName)
	if err == nil {
		notify, stop := make(chan struct{}), make(chan struct{})
		go func() {
			listening(notify, stop)
			can.Disconnect()
		}()
		f(can)
		close(stop)
		<-notify
	}
	return err
}

func listening(notify chan struct{}, stop <-chan struct{}) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGTSTP)
	select {
	case <-stop:
	case <-quit:
	}
	close(notify)
}
