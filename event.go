/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import "time"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	TimerId      uint        // Unique ID for each timer created
	EventHandler func(Event) // Handler for an emitted event
	EventNS      uint        // Event namespace
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

	Emit(Event)                      // Emit an event on the bus
	NewHandler(string, EventHandler) // Register an event handler for an event name
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
