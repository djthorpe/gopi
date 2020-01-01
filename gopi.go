/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

////////////////////////////////////////////////////////////////////////////////
// TYPES

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

// Abstract logging interface
type Logger interface {
	Unit

	Name() string              // Return unit name
	Error(error) error         // Output logging messages
	Debug(args ...interface{}) // Debug output
	IsDebug() bool             // Return IsDebug flag
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
