/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"context"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EventNS      uint                         // Event namespace
	EventHandler func(Event, context.Context) // Handler for an emitted event
	TimerId      uint                         // Unique ID for each timer created
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Timer unit - sends out messages on the event bus
type Timer interface {
	Unit

	NewTicker(time.Duration) TimerId // Create periodic event at interval
	NewTimer(time.Duration) TimerId  // Create one-shot event after interval
	Cancel(TimerId) error            // Cancel events
}

// Bus unit - handles events
type Bus interface {
	Unit

	Emit(Event)                                                    // Emit an event on the bus
	NewHandler(string, EventNS, EventHandler, time.Duration) error // Register an event handler for an event name
	DefaultHandler(EventNS, EventHandler, time.Duration) error     // Register default handler for events
}

// Event emitted on the event bus
type Event interface {
	Source() Unit       // Source of the event
	Name() string       // Name of the event
	NS() EventNS        // Namespace for the event
	Value() interface{} // Any value associated with the event
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EVENT_NS_DEFAULT EventNS = iota // Default event namespace
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v EventNS) String() string {
	switch v {
	case EVENT_NS_DEFAULT:
		return "EVENT_NS_DEFAULT"
	default:
		return "[?? Invalid EventNS value]"
	}
}
