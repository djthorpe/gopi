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

type Subscriber int

type Event interface {
	// Source of the event
	Source() Driver

	// Generic name of the event
	Name() string
}

type Publisher interface {
	// Subscribe to events emitted. Returns unique subscriber
	// identifier and channel on which events are emitted
	Subscribe() (Subscriber, chan Event)

	// Unsubscribe from events emitted
	Unsubscribe(Subscriber)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// SUBSCRIBER_NONE is returned when a subscribe call is not implemented
	SUBSCRIBER_NONE Subscriber = 0
)
