/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package timer

import (
	"fmt"

	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type timerEvent struct {
	source  gopi.Timer
	eventId gopi.EventId
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func newTimerEvent(source gopi.Timer, eventId gopi.EventId) gopi.Event {
	return &timerEvent{source, eventId}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Event

func (*timerEvent) Name() string {
	return "gopi.TimerEvent"
}

func (*timerEvent) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

func (this *timerEvent) Source() gopi.Unit {
	return this.source
}

func (this *timerEvent) EventId() gopi.EventId {
	return this.eventId
}

func (this *timerEvent) Value() interface{} {
	return this.eventId
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *timerEvent) String() string {
	return fmt.Sprintf("<%v ns=%v value=%v>", this.Name(), this.NS(), this.Value())
}
