/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"context"
	"io"
	"os"
	"testing"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// MainCommandFunc is the main handler for command line tool
	MainCommandFunc func(App, []string) error

	// MainTestFunc is the main handler for a testing tool
	MainTestFunc func(App, *testing.T)

	// FlagNS is the namespace for a flag
	FlagNS uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// App encapsulates the lifecycle of a running application
// which is created through the 'app' package and run by calling
// the Run method. A running application has access to unit instances
// (for example, logger, timer and event bus).
type App interface {

	// Run the application and return the return code
	Run() int

	// WaitForSignal stalls the application until a signal is returned,
	// context is cancelled or deadline exceeded
	WaitForSignal(context.Context, ...os.Signal) error // Wait for interrupt signal with context

	// Flags returns the set of key/value configuration parameters
	Flags() Flags

	// Emit events
	Emit(Event)

	// Return unit instances
	UnitInstance(string) Unit // Return singular unit for name
	Log() Logger              // Return logger unit
	Timer() Timer             // Return timer unit
	Bus() Bus                 // Return event bus unit
	Platform() Platform       // Return hardware platform unit
	Display() Display         // Return display
	Surfaces() SurfaceManager // Return surface manager
	Input() InputManager      // Return input manager
	Fonts() FontManager       // Return font manager
	GPIO() GPIO               // Return GPIO interface
	I2C() I2C                 // Return I2C interface
	SPI() SPI                 // Return SPI interface
	LIRC() LIRC               // Return LIRC interface
}

// Flags encapsulates a set of key/value pairs in several namespaces
// with parsing of command-line flags in the default namespace
type Flags interface {
	Name() string                // Return name of tool
	Parse([]string) error        // Parse command-line flags
	Parsed() bool                // Returns true if Parsed() has been called
	Args() []string              // Return command-line arguments
	HasFlag(string, FlagNS) bool // HasFlag returns true if a flag exists

	Usage(io.Writer)   // Write out usage for the application
	Version(io.Writer) // Write out version for the application
	//	SetUsage(func(io.Writer))    // Set command usage function

	// Define flags in default namespace
	FlagBool(name string, value bool, usage string) *bool
	FlagString(name, value, usage string) *string
	FlagDuration(name string, value time.Duration, usage string) *time.Duration
	FlagInt(name string, value int, usage string) *int
	FlagUint(name string, value uint, usage string) *uint
	FlagFloat64(name string, value float64, usage string) *float64

	// Get flag values
	GetBool(string, FlagNS) bool
	GetString(string, FlagNS) string
	GetDuration(string, FlagNS) time.Duration
	GetInt(string, FlagNS) int
	GetUint(string, FlagNS) uint
	GetFloat64(string, FlagNS) float64

	// Set flag values
	SetBool(string, FlagNS, bool)
	SetString(string, FlagNS, string)
	SetDuration(string, FlagNS, time.Duration)
	SetInt(string, FlagNS, int)
	SetUint(string, FlagNS, uint)
	SetFloat64(string, FlagNS, float64)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FLAG_NS_DEFAULT FlagNS = iota
	FLAG_NS_VERSION
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v FlagNS) String() string {
	switch v {
	case FLAG_NS_DEFAULT:
		return "FLAG_NS_DEFAULT"
	case FLAG_NS_VERSION:
		return "FLAG_NS_VERSION"
	default:
		return "[?? Invalid FlagNS value]"
	}
}
