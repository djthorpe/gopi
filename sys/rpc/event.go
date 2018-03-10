/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpc

import (
	"fmt"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Event is the RPC event
type Event struct {
	source gopi.Driver
	t      gopi.RPCEventType
	r      *gopi.RPCServiceRecord
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return the type of event
func (e *Event) Type() gopi.RPCEventType {
	return e.t
}

// Return the service record
func (e *Event) ServiceRecord() *gopi.RPCServiceRecord {
	return e.r
}

// Return name of event
func (*Event) Name() string {
	return "RPCEvent"
}

// Return source of event
func (e *Event) Source() gopi.Driver {
	return e.source
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e *Event) String() string {
	if e.r != nil {
		return fmt.Sprintf("<rpc.Event>{ type=%v record=%v }", e.Type(), e.ServiceRecord())
	} else {
		return fmt.Sprintf("<rpc.Event>{ type=%v }", e.Type())
	}
}
