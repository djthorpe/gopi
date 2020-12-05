package mdns

import (
	"fmt"

	"github.com/djthorpe/gopi/v3"
	"github.com/miekg/dns"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type msgevent struct {
	*dns.Msg
}

///////////////////////////////////////////////////////////////////////////////
// NEW

func NewMsgEvent(msg *dns.Msg) gopi.Event {
	return &msgevent{msg}
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
	return str + ">"
}
