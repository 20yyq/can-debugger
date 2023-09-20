// @@
// @ Author       : Eacher
// @ Date         : 2023-09-15 14:37:58
// @ LastEditTime : 2023-09-19 08:12:30
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

	"golang.org/x/sys/unix"
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
	types		string
)

type Interface struct {
	iface	*net.Interface
	conn 	*netlink.NetlinkRoute
}

func InitFlagArge(f *flag.FlagSet) {
	f.UintVar(&bitrate, "bitrate", 125000, "bitrate usage uint")
	f.StringVar(&types, "type", "can", "type usage string")
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
			if err = iface.Up(); err != nil {
				err = iface.Up()
			}
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
	sm.Attrs = append(sm.Attrs, packet.IfInfomsg{Flags: syscall.IFF_UP, Change: syscall.IFF_UP, Index: int32(ifi.iface.Index)})
	rm := netlink.ReceiveNLMessage{Data: make([]byte, 1024)}
	err := ifi.conn.Exchange(&sm, &rm)
	if err == nil {
		if rm.MsgList[0].Header.Type != syscall.RTM_NEWLINK {
			err = DeserializeNlMsgerr(rm.MsgList[0])
		} else {
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
	sm.Attrs = append(sm.Attrs, packet.IfInfomsg{Change: syscall.IFF_UP, Index: int32(ifi.iface.Index)})
	rm := netlink.ReceiveNLMessage{Data: make([]byte, 1024)}
	err := ifi.conn.Exchange(&sm, &rm)
	if err == nil {
		if rm.MsgList[0].Header.Type != syscall.RTM_NEWLINK {
			err = DeserializeNlMsgerr(rm.MsgList[0])
		} else {
			ifi.iface.Flags &= 0xFFFFFFFE
		}
	}
	return err
}

func (ifi *Interface) SetBitrate() error {
	sm := netlink.SendNLMessage{
		NlMsghdr: &packet.NlMsghdr{Type: syscall.RTM_NEWLINK, Flags: syscall.NLM_F_REQUEST|syscall.NLM_F_ACK, Seq: randReq()},
	}
	sm.Attrs = append(sm.Attrs, packet.IfInfomsg{Index: int32(ifi.iface.Index)})
	load := packet.NlAttr{&syscall.NlAttr{uint16(len(types) + packet.SizeofNlAttr), unix.IFLA_INFO_KIND}, append([]byte(types), 0x00)}.WireFormat()
	data := (packet.CANBitTiming{Bitrate: uint32(bitrate)}).WireFormat()
	data = packet.NlAttr{&syscall.NlAttr{uint16(len(data) + packet.SizeofNlAttr), unix.IFLA_CAN_BITTIMING}, data}.WireFormat()
	load = append(load, packet.NlAttr{&syscall.NlAttr{uint16(len(data) + packet.SizeofNlAttr), unix.IFLA_INFO_DATA}, data}.WireFormat()...)
	sm.Attrs = append(sm.Attrs, packet.NlAttr{&syscall.NlAttr{uint16(len(load) + packet.SizeofNlAttr), syscall.IFLA_LINKINFO}, load})
	rm := netlink.ReceiveNLMessage{Data: make([]byte, 1024)}
	err := ifi.conn.Exchange(&sm, &rm)
	if err == nil {
		if rm.MsgList[0].Header.Type != syscall.RTM_NEWLINK {
			err = DeserializeNlMsgerr(rm.MsgList[0])
		}
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
	if len(nlm.Data) >= packet.SizeofNlMsgerr {
		msg := packet.NewNlMsgerr(([packet.SizeofNlMsgerr]byte)(nlm.Data[:packet.SizeofNlMsgerr]))
		if msg.Error < 0 {
			msg.Error *= -1
		}
		if msg.Error > 0 {
			return syscall.Errno(msg.Error)
		}
	}
	return nil
}
