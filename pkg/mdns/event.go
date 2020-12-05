package mdns

import (
	"fmt"
	"net"

	"github.com/djthorpe/gopi/v3"
	"github.com/miekg/dns"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type msgevent struct {
	*dns.Msg
	net.Addr
	ifIndex int
}

///////////////////////////////////////////////////////////////////////////////
// NEW

func NewMsgEvent(msg *dns.Msg, addr net.Addr, ifIndex int) gopi.Event {
	return &msgevent{msg, addr, ifIndex}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC PROPERTIES

func (this *msgevent) Name() string {
	if this.Msg == nil {
		return ""
	} else {
		return this.Msg.MsgHdr.String()
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *msgevent) String() string {
	str := "<dns.msg"
	str += fmt.Sprintf(" name=%q", this.Name())
	if this.Addr != nil {
		str += fmt.Sprintf(" addr=%v", this.Addr)
	}
	if this.ifIndex >= 0 {
		str += fmt.Sprintf(" ifIndex=%v", this.ifIndex)
	}
	return str + ">"
}
