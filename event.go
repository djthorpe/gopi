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

	// EventHandlerFunc is the handler for an emitted event, which
	// should cancel when context.Done() signal is received
	EventHandlerFunc func(context.Context, App, Event)

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
	NewHandler(EventHandler) error

	// DefaultHandler registers a default handler for an event namespace
	DefaultHandler(EventNS, EventHandlerFunc) error
}

// Event emitted on the event bus
type Event interface {
	Source() Unit       // Source of the event
	Name() string       // Name of the event
	NS() EventNS        // Namespace for the event
	Value() interface{} // Any value associated with the event
}

// EventHandler defines how an emitted event is handled in the application
type EventHandler struct {
	// The name of the event
	Name string

	// The handler function for the event
	Handler EventHandlerFunc

	// The namespace of the event, usually 0
	EventNS EventNS

	// The timeout value for the handler, usually 0
	Timeout time.Duration
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EVENT_NS_DEFAULT EventNS = iota // Default event namespace
	EVENT_NS_MAX             = EVENT_NS_DEFAULT
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	NullEvent Event
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
