package gopi

import (
	"context"
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Config interface {
	// Parse
	Parse() error

	// Return arguments
	Args() []string

	// Define flags
	FlagString(string, string, string) *string
	FlagBool(string, bool, string) *bool
	FlagUint(string, uint, string) *uint

	// Get config values
	GetString(string) string
	GetBool(string) bool
	GetUint(string) uint
}

// Unit marks an singleton object
type Unit struct{}

/////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

func (this *Unit) Define(Config) error       { /* NOOP */ return nil }
func (this *Unit) New(Config) error          { /* NOOP */ return nil }
func (this *Unit) Run(context.Context) error { /* NOOP */ return nil }
func (this *Unit) Dispose() error            { /* NOOP */ return nil }
