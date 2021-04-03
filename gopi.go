package gopi

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Config interface {
	Parse() error     // Parse command line arguments
	Args() []string   // Return arguments, not including flags
	Usage(string)     // Print out usage for all or specific command
	Version() Version // Return version information

	// Define flags
	FlagString(string, string, string, ...string) *string
	FlagBool(string, bool, string, ...string) *bool
	FlagUint(string, uint, string, ...string) *uint
	FlagInt(string, int, string, ...string) *int
	FlagDuration(string, time.Duration, string, ...string) *time.Duration
	FlagFloat(string, float64, string, ...string) *float64

	// Define a command with name, description, calling function
	Command(string, string, CommandFunc) error

	// Get command from provided arguments
	GetCommand([]string) (Command, error)

	// Get flag values
	GetString(string) string
	GetBool(string) bool
	GetUint(string) uint
	GetInt(string) int
	GetDuration(string) time.Duration
	GetFloat(string) float64
}

// CommandFunc is the function signature for running a command
type CommandFunc func(context.Context) error

// Command is determined from parsed arguments
type Command interface {
	Name() string              // Return command name
	Usage() (string, string)   // Return command syntax and description
	Args() []string            // Return command arguments
	Run(context.Context) error // Run the command
}

type Version interface {
	Name() string                      // Return process name
	Version() (string, string, string) // Return tag, branch and hash
	BuildTime() time.Time              // Return time of process compilation
	GoVersion() string                 // Return go compiler version
}

// Logger outputs information and debug messages
type Logger interface {
	Print(args ...interface{})              // Output logging
	Debug(args ...interface{})              // Output debugging information
	Printf(fmt string, args ...interface{}) // Output logging with format
	Debugf(fmt string, args ...interface{}) // Output debugging with format
	IsDebug() bool                          // IsDebug returns true if debug flag is set
	T() *testing.T                          // When testing, provides testing context
}

// Event is an emitted event
type Event interface {
	Name() string // Return name of the event
}

// Publisher emits events and allows for subscribing to emitted events
type Publisher interface {
	// Emit an event, which can block if second argument is true
	Emit(Event, bool) error

	// Subscribe to events
	Subscribe() <-chan Event

	// Unsubscribe from events
	Unsubscribe(<-chan Event)
}

// Promises runs chains of events in the background
type Promises interface {
	// Create a promise with a function
	Do(context.Context, func(context.Context, interface{}) (interface{}, error), interface{}) Promise
}

// Promise is run in a chain, and finally calls given function
type Promise interface {
	// Chain a function to a promise
	Then(func(context.Context, interface{}) (interface{}, error)) Promise

	// Finally runs the promise in the background and optionally waits for it to complete
	// then returns any error if not running in background
	Finally(func(interface{}, error) error, bool) error
}

/////////////////////////////////////////////////////////////////////
// UNITS

// Unit marks an singleton object
type Unit struct{}

/////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

func (this *Unit) Define(Config) error       { /* NOOP */ return nil }
func (this *Unit) New(Config) error          { /* NOOP */ return nil }
func (this *Unit) Run(context.Context) error { /* NOOP */ return nil }
func (this *Unit) Dispose() error            { /* NOOP */ return nil }

/////////////////////////////////////////////////////////////////////
// REQUIRE

// Call Require with a set of values and if any of them are nil then panic
func (this *Unit) Require(units ...interface{}) {
	for _, v := range units {
		if v == nil {
			fmt.Fprintf(os.Stderr, "Require: %v not satisfied\n", reflect.TypeOf(v))
			os.Exit(-1)
		}
	}
}
