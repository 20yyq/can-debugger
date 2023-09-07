// @@
// @ Author       : Eacher
// @ Date         : 2023-09-06 14:47:15
// @ LastEditTime : 2023-09-07 14:10:35
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  : 
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/read/can.go
// @@
package read

import (
	"fmt"
	
	"github.com/20yyq/can/sockcan"
	"github.com/20yyq/packet/can"
)

func Run(c *sockcan.Can) {
	fmt.Println("Running...")
	for {
		if frame, err := c.ReadFrame(); err == nil {
			go printFrame(frame)
			continue
		}
		break
	}
}

func printFrame(f can.Frame) {
	fmt.Println("frame", f)
}
