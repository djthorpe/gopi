/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

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

// Publisher interface for Subscribe/Emit mechanism for arbitary messages
type Publisher interface {
	// Subscribe to a queue with a capacity
	Subscribe(queue uint, capacity int) <-chan interface{}

	// Unsubscribe channel from queue, returning true if successful
	Unsubscribe(<-chan interface{}) bool
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
