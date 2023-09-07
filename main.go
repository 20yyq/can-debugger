// @@
// @ Author       : Eacher
// @ Date         : 2023-09-02 11:02:03
// @ LastEditTime : 2023-09-07 14:10:17
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  : 
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/main.go
// @@
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/20yyq/can/flag"
	"github.com/20yyq/can/sockcan"
	"github.com/20yyq/can/read"
	"github.com/20yyq/can/write"
)

func main() {
	can, err := sockcan.NewCan()
	if err != nil {
		fmt.Println("can start err: ", err)
		return
	}
	stop := make(chan struct{})
	go listening(stop, can)
	switch flag.DebuggerName() {
	case "read":
		read.Run(can)
	case "write":
		if err = flag.Init(write.InitFlagArge); err == nil {
			write.Run(can)
		}
	case "info":
		read.Run(can)
	default:
		fmt.Println("can start default")
	}
	close(stop)
}

func listening(stop chan struct{}, c *sockcan.Can) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	select{
	case <-stop:
	case <-quit:
		close(stop)
	}
	c.Disconnect()
	fmt.Println("End")
	os.Exit(0)
}