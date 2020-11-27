<<<<<<< HEAD
=======
/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

>>>>>>> master
package gopi

import (
	"context"
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
	FlagString(string, string, string) *string
	FlagBool(string, bool, string) *bool
	FlagUint(string, uint, string) *uint
	FlagDuration(string, time.Duration, string) *time.Duration

	// Define commands
	Command(string, string, CommandFunc) error // Append a command with name and usage arguments

	// Get values
	GetString(string) string
	GetBool(string) bool
	GetUint(string) uint
	GetDuration(string) time.Duration
	GetCommand([]string) Command // Get command from provided arguments
}

// CommandFunc is the function signature for running a command
type CommandFunc func(context.Context) error

// Command is determined from parsed arguments
type Command interface {
	Name() string              // Return command name
	Usage() string             // Return usage information
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
