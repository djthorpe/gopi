package mdns

import (
	"fmt"

	"github.com/djthorpe/gopi/v3"
	"github.com/miekg/dns"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type dnsevent struct {
	msg *dns.Msg
}

///////////////////////////////////////////////////////////////////////////////
// NEW

func NewDNSEvent(msg *dns.Msg) gopi.Event {
	this := new(dnsevent)
	this.msg = msg
	return this
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC PROPERTIES

func (this *dnsevent) Name() string {
	if this.msg == nil {
		return ""
	} else {
		return this.msg.MsgHdr.String()
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *dnsevent) String() string {
	str := "<dnsevent"
	str += fmt.Sprintf(" name=%q", this.Name())
	return str + ">"
}
