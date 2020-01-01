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
	// EventNS is the namespace in which events are emitted, usually
	// EVENT_NS_DEFAULT
	EventNS uint

	// EventHandler is the handler for an emitted event, which should cancel
	// when context.Done() signal is received
	EventHandler func(context.Context, Event)

	// EventId is a unique ID for each event
	EventId uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Timer unit - sends out messages on the event bus
type Timer interface {
	Unit

	NewTicker(time.Duration) EventId // Create periodic event at interval
	NewTimer(time.Duration) EventId  // Create one-shot event after interval
	Cancel(EventId) error            // Cancel events
}

// Bus unit - handles events
type Bus interface {
	Unit

	// Emit an event on the bus
	Emit(Event)

	// NewHandler registers an event handler for an event name
	NewHandler(string, EventHandler) error

	// NewHandlerEx registers an event handler for an event name, in a
	// particular event namespace, and with a timeout value for the handler
	NewHandlerEx(string, EventNS, EventHandler, time.Duration) error

	// DefaultHandler registers a default handler for an event namespace
	DefaultHandler(EventNS, EventHandler) error
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
	EVENT_NS_MAX             = EVENT_NS_DEFAULT
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
