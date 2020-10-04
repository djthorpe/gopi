package gopi

import (
	"context"
)

/////////////////////////////////////////////////////////////////////
// INTERFACES

type Config interface {
	// Parse
	Parse() error

	// Define flags
	String(string, string, string) *string
	Bool(string, bool, string) *bool

	// Get configuration values
	GetString(string) string
	GetBool(string) bool
}

// Unit marks an singleton object
type Unit struct{}

/////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

func (this *Unit) Define(Config) error       { /* NOOP */ return nil }
func (this *Unit) New(Config) error          { /* NOOP */ return nil }
func (this *Unit) Run(context.Context) error { /* NOOP */ return nil }
func (this *Unit) Dispose() error            { /* NOOP */ return nil }
