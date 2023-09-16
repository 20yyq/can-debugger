// @@
// @ Author       : Eacher
// @ Date         : 2023-09-02 11:02:03
// @ LastEditTime : 2023-09-16 16:00:19
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  : 
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/main.go
// @@
package main

import (
	"os"
	"fmt"

	"github.com/20yyq/can/flag"
)

func main() {
	err := flag.Init()
	if err == nil {
		if err = flag.Runing(); err == nil {
			fmt.Println("End")
			os.Exit(0)
		}
	}
	fmt.Println("can runing err: ", err)
	os.Exit(1)
}
