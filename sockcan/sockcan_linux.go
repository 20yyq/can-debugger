// @@
// @ Author       : Eacher
// @ Date         : 2023-09-06 14:55:23
// @ LastEditTime : 2023-09-07 14:12:38
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  : 
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/sockcan/sockcan_linux.go
// @@
package sockcan

import (
	"net"
	"os"
	"io"
	"syscall"

	"github.com/20yyq/can/flag"
	"github.com/20yyq/packet/can"
	
	"golang.org/x/sys/unix"
)

func NewCan() (*Can, error) {
	iface, err := net.InterfaceByName(flag.InterfaceName())
	if err == nil {
		fd, _ := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, unix.CAN_RAW)
		if err = unix.Bind(fd, &unix.SockaddrCAN{Ifindex: iface.Index}); err == nil {
			return &Can{rwc: os.NewFile(uintptr(fd), flag.InterfaceName())}, nil
		}
	}
	return nil, err
}

type HandlerFunc func(can.Frame)

type Can struct {
	rwc 	io.ReadWriteCloser
}

func (c *Can) ReadFrame() (f can.Frame, err error) {
	var b [can.FrameLength]byte
	_, err = c.rwc.Read(b[:])
	if err == nil {
		f = can.NewFrame(b)
	} else if err == io.EOF {
		c.rwc.Close()
	}
	return f, err
}

func (c *Can) WriteFrame(frame can.Frame) error {
	_, err := c.rwc.Write(frame.WireFormat())
	return err
}

func (c *Can) Disconnect() error {
	return c.rwc.Close()
}