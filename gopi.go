package gopi

import (
	"context"
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Config interface {
	Parse() error   // Parse command line arguments
	Args() []string // Return arguments, not including flags
	Usage(string)   // Print out usage for all or specific command

	// Define flags
	FlagString(string, string, string) *string
	FlagBool(string, bool, string) *bool
	FlagUint(string, uint, string) *uint

	// Define commands
	Command(string, string, CommandFunc) error // Append a command with name and usage arguments

	// Get values
	GetString(string) string
	GetBool(string) bool
	GetUint(string) uint
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

// Unit marks an singleton object
type Unit struct{}

// Event is a generic emitted event
type Event interface{}

// Logger outputs information and debug messages
type Logger interface {
	Print(args ...interface{}) // Output logging
	Debug(args ...interface{}) // Output debugging information
	IsDebug() bool             // IsDebug returns true if debug flag is set
}

/////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

func (this *Unit) Define(Config) error       { /* NOOP */ return nil }
func (this *Unit) New(Config) error          { /* NOOP */ return nil }
func (this *Unit) Run(context.Context) error { /* NOOP */ return nil }
func (this *Unit) Dispose() error            { /* NOOP */ return nil }
