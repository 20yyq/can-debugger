// @@
// @ Author       : Eacher
// @ Date         : 2023-09-06 14:47:15
// @ LastEditTime : 2024-01-09 16:29:16
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  :
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/read/can.go
// @@
package read

import (
	"fmt"

	"github.com/20yyq/can-debugger/sockcan"
	"github.com/20yyq/packet/can"
)

func Run(c *sockcan.Can) {
	fmt.Println("Running...")
	for {
		frame, err := c.ReadFrame()
		if err == nil {
			go printFrame(frame)
			continue
		}
		fmt.Println("ReadFrame err", err)
		break
	}
}

func printFrame(f can.Frame) {
	fmt.Println("frame", f)
}
