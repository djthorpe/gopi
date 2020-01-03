/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	source  gopi.Unit
	type_   gopi.RPCEventType
	service gopi.RPCServiceRecord
	ttl     time.Duration
}

func NewEvent(source gopi.Unit, type_ gopi.RPCEventType, service gopi.RPCServiceRecord, ttl time.Duration) gopi.RPCEvent {
	return &event{source, type_, service, ttl}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Event

func (*event) Name() string { return "gopi.RPCEvent" }

func (this *event) Source() gopi.Unit {
	return this.source
}

func (this *event) Value() interface{} {
	return this.service
}

func (this *event) NS() gopi.EventNS {
	return gopi.EVENT_NS_DEFAULT
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.RPCEvent

func (this *event) Type() gopi.RPCEventType {
	return this.type_
}

func (this *event) Service() gopi.RPCServiceRecord {
	return this.service
}

func (this *event) TTL() time.Duration {
	return this.ttl
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<" + this.Name() + " type=" + fmt.Sprint(this.type_)
	if this.service.Name != "" {
		str += " service=" + fmt.Sprint(this.service)
		str += " ttl=" + fmt.Sprint(this.ttl)
	}
	return str + ">"
}
