// @@
// @ Author       : Eacher
// @ Date         : 2023-09-15 14:37:58
// @ LastEditTime : 2023-09-16 15:59:42
// @ LastEditors  : Eacher
// @ --------------------------------------------------------------------------------<
// @ Description  : 
// @ --------------------------------------------------------------------------------<
// @ FilePath     : /20yyq/can-debugger/iface/iface.go
// @@
package iface

import (
	"fmt"
	"flag"
	"net"
	"time"
	"syscall"

	"github.com/20yyq/netlink"
	"github.com/20yyq/packet"
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
	bitrate		uint
	up			bool
)

type Interface struct {
	iface	*net.Interface
	conn 	*netlink.NetlinkRoute
}

func InitFlagArge(f *flag.FlagSet) {
	f.UintVar(&bitrate, "bitrate", 1250000, "bitrate usage uint")
	f.BoolVar(&up, "up", false, "up usage bool")
	f.Usage = help
}

func Run(dev string) error {
	fmt.Println("Running...")
	iface, err := newInterface(dev)
	if err == nil {
		if !up {
			return iface.Down()
		}
		if err = iface.SetBitrate(); err == nil {
			err = iface.Up()
		}
		iface.conn.Close()
	}
	return err
}

func newInterface(dev string) (*Interface, error) {
	iface := &Interface{}
	var err error
	if iface.iface, err = net.InterfaceByName(dev); err != nil {
		return nil, err
	}
	iface.conn = &netlink.NetlinkRoute{
		DevName: iface.iface.Name,
		Sal: &syscall.SockaddrNetlink{Family: syscall.AF_NETLINK},
	}
	return iface, iface.conn.Init() 
}

func (ifi *Interface) Up() error {
	if ifi.iface.Flags&0x01 != 0 {
		return nil
	}
	sm := netlink.SendNLMessage{
		NlMsghdr: &packet.NlMsghdr{Type: syscall.RTM_NEWLINK, Flags: syscall.NLM_F_REQUEST|syscall.NLM_F_ACK, Seq: randReq()},
	}
	sm.Attrs = append(sm.Attrs, packet.IfInfomsg{Family: syscall.AF_UNSPEC, Flags: syscall.IFF_UP, Change: syscall.IFF_UP, Index: int32(ifi.iface.Index)})
	rm := netlink.ReceiveNLMessage{Data: make([]byte, 128)}
	err := ifi.conn.Exchange(&sm, &rm)
	if err == nil {
		if err = DeserializeNlMsgerr(rm.MsgList[0]); err == nil {
			ifi.iface.Flags |= 0x01
		}
	}
	return err
}

func (ifi *Interface) Down() error {
	if ifi.iface.Flags&0x01 != 1 {
		return nil
	}
	sm := netlink.SendNLMessage{
		NlMsghdr: &packet.NlMsghdr{Type: syscall.RTM_NEWLINK, Flags: syscall.NLM_F_REQUEST|syscall.NLM_F_ACK, Seq: randReq()},
	}
	sm.Attrs = append(sm.Attrs, packet.IfInfomsg{Family: syscall.AF_UNSPEC, Change: syscall.IFF_UP, Index: int32(ifi.iface.Index)})
	rm := netlink.ReceiveNLMessage{Data: make([]byte, 128)}
	err := ifi.conn.Exchange(&sm, &rm)
	if err == nil {
		if err = DeserializeNlMsgerr(rm.MsgList[0]); err == nil {
			ifi.iface.Flags &= 0xFFFFFFFE
		}
	}
	return err
}

func (ifi *Interface) SetBitrate() error {
	sm := netlink.SendNLMessage{
		NlMsghdr: &packet.NlMsghdr{Type: syscall.RTM_NEWLINK, Flags: syscall.NLM_F_REQUEST|syscall.NLM_F_ACK, Seq: randReq()},
	}
	sm.Attrs = append(sm.Attrs, packet.IfInfomsg{Family: syscall.AF_UNSPEC, Index: int32(ifi.iface.Index)})
	timing := (packet.CANBitTiming{Bitrate: uint32(bitrate)}).WireFormat()
	sm.Attrs = append(sm.Attrs, packet.NlAttr{&syscall.NlAttr{uint16(len(timing) + packet.SizeofNlAttr), syscall.IFLA_LINKINFO}, timing})
	rm := netlink.ReceiveNLMessage{Data: make([]byte, 256)}
	err := ifi.conn.Exchange(&sm, &rm)
	if err == nil {
		err = DeserializeNlMsgerr(rm.MsgList[0])
	}
	return err
}

func help() {
	fmt.Print(helpOutput)
}

func randReq() uint32 {
	return uint32(time.Now().UnixNano() & 0xFFFFFFFF)
}

func DeserializeNlMsgerr(nlm *packet.NetlinkMessage) error {
	if len(nlm.Data) < packet.SizeofNlMsgerr {
		return syscall.Errno(34)
	}
	msg := packet.NewNlMsgerr(([packet.SizeofNlMsgerr]byte)(nlm.Data[:packet.SizeofNlMsgerr]))
	if msg.Error < 0 {
		msg.Error *= -1
	}
	if msg.Error > 0 {
		return syscall.Errno(msg.Error)
	}
	return nil
}
