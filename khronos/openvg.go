package khronos

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Store state for the non-abstract input driver
type State struct {
	driver Driver
}

// Abstract configuration which is used to open and return the
// concrete driver
type Config interface {
	// Opens the driver from configuration, or returns error
	Open() (Driver, error)
}

// Abstract driver interface
type Driver interface {
	// Close closes the driver and frees the underlying resources
	Close() error
}

////////////////////////////////////////////////////////////////////////////////
// Opener interface

// Open opens a connection to the touchscreen with the given driver.
func Open2(config Config) (Driver, error) {
	driver, err := config.Open()
	if err != nil {
		return nil, err
	}
	return &State{driver}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Driver interface

// Provides human-readable version
func (state *State) String() string {
	return fmt.Sprintf("<OpenVG>{%v}", state.driver)
}

// Closes the device and frees the resources
func (state *State) Close() error {
	return state.driver.Close()
}


