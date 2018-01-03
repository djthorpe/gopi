/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import "time"

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Publisher is an interface for drivers which accept subscription
// and unsubscription requests
type Publisher interface {
	// Subscribe to events emitted. Returns channel on which events
	// are emitted or nil if this driver does not implement events
	Subscribe() <-chan Event

	// Unsubscribe from events emitted
	Unsubscribe(<-chan Event)

	// Emit an event to subscribers
	Emit(Event)
}

// Event is a generic event which is emitted through a channel
type Event interface {
	// Source of the event
	Source() Driver

	// Name of the event
	Name() string
}

// GPIOEvent implements an event from the GPIO driver
type GPIOEvent interface {
	Event

	// Pin returns the pin on which the event occurred
	Pin() GPIOPin

	// Edge returns whether the pin value is rising or falling
	// or will return NONE if not defined
	Edge() GPIOEdge
}

// LIRCEvent implements an event from the LIRC driver
type LIRCEvent interface {
	Event

	// The type of message
	Type() LIRCType

	// The value
	Value() uint32
}

// TimerEvent is emitted by the timer driver on maturity
type TimerEvent interface {
	Event

	Timestamp() time.Time
	UserInfo() interface{}
}
