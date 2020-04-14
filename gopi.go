/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// Channel is an arbitary communication channel
	Channel uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Unit configuration interface
type Config interface {
	Name() string             // Returns name of the unit
	New(Logger) (Unit, error) // Opens the driver from configuration, or returns error
}

// Unit interface
type Unit interface {
	Close() error   // Close closes the driver and frees the underlying resources
	String() string // String returns a string representation of the unit
}

// Publisher interface for Emit/Receive mechanism for arbitary messages
type Publisher interface {
	// Emit sends values to be received by handlers
	Emit(queue uint, value interface{})

	// Subscribe and handle messages which are emitted
	Subscribe(queue uint, callback func(value interface{})) Channel

	// Unsubscribe from a channel
	Unsubscribe(Channel)
}

// Publisher/Subscriber interface for arbitary messages
type PubSub interface {
	Unit

	// Emit sends values to be received by handlers
	Emit(value interface{})

	// Subscribe and handle messages which are emitted
	Subscribe() <-chan interface{}

	// Unsubscribe from a channel
	Unsubscribe(<-chan interface{})
}

// Abstract logging interface
type Logger interface {
	Unit

	// Name returns the name of the logger
	Name() string

	// Clone returns a new logger with a different name
	Clone(string) Logger

	// Error logs an error, and returns an error with the name prefixed
	Error(error) error

	// Warn outputs informational message with warning severity
	Warn(args ...interface{})

	// Info outputs informational message
	Info(args ...interface{})

	// Debug will log a debug message when debugging is on
	Debug(args ...interface{})

	// IsDebug returns true if debugging is enabled
	IsDebug() bool
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Config

func New(config Config, log Logger) (Unit, error) {
	if driver, err := config.New(log); err != nil {
		return nil, err
	} else {
		return driver, nil
	}
}
