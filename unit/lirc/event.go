/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package lirc

import (
	// Frameworks
	"fmt"

	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type lirc_event struct {
	source gopi.Unit
	value  uint32
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.LIRCEvent

func NewEvent(source gopi.Unit, value uint32) gopi.Event {
	return &lirc_event{source, value}
}

func (this *lirc_event) Name() string {
	return "gopi.LIRCEvent"
}

func (this *lirc_event) Source() gopi.Unit {
	return this.source
}

func (this *lirc_event) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *lirc_event) Type() gopi.LIRCType {
	return gopi.LIRCType(this.value & 0xFF000000)
}

func (this *lirc_event) Value() interface{} {
	return this.value & 0x00FFFFFF
}

func (this *lirc_event) String() string {
	return "<lirc.event" +
		" type=" + fmt.Sprint(this.Type()) +
		" value=" + fmt.Sprint(this.Value()) +
		">"
}
