/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"fmt"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Event is the RPC event
type event struct {
	s gopi.Driver
	t gopi.RPCEventType
	r *gopi.RPCServiceRecord
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewEvent(source gopi.Driver, event_type gopi.RPCEventType, service_record *gopi.RPCServiceRecord) *event {
	return &event{s: source, t: event_type, r: service_record}
}

// Return the type of event
func (this *event) Type() gopi.RPCEventType {
	return this.t
}

// Return the service record
func (this *event) ServiceRecord() *gopi.RPCServiceRecord {
	return this.r
}

// Return name of event
func (*event) Name() string {
	return "RPCEvent"
}

// Return source of event
func (this *event) Source() gopi.Driver {
	return this.s
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	if this.r != nil {
		return fmt.Sprintf("<rpc.event>{ type=%v record=%v }", this.t, this.r)
	} else {
		return fmt.Sprintf("<rpc.event>{ type=%v }", this.t)
	}
}
