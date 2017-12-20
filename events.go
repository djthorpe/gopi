/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Event is a generic event which is emitted through a channel
type Event interface {
	// Source of the event
	Source() Driver

	// Name of the event
	Name() string
}

type Publisher interface {
	// Subscribe to events emitted. Returns channel on which events
	// are emitted or nil if this driver does not implement events
	Subscribe() chan Event

	// Unsubscribe from events emitted
	Unsubscribe(chan Event)
}
