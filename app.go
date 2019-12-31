/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MainCommandFunc func(App, []string) error
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type App interface {
	Run() int // Run application

	Flags() Flags // Return command-line flags
	Log() Logger  // Return logger unit
	Timer() Timer // Return timer unit
	Bus() Bus     // Return event bus unit

	Unit(string) Unit    // Return singular unit for name
	Units(string) []Unit // Return multiple units for name
}

type Flags interface {
	Name() string        // Return name of flagset
	Args() []string      // Args returns the command line arguments
	HasFlag(string) bool // HasFlag returns true if a flag exists

	FlagBool(name string, value bool, usage string) *bool
	FlagString(name, value, usage string) *string
	FlagDuration(name string, value time.Duration, usage string) *time.Duration
	FlagInt(name string, value int, usage string) *int
	FlagUint(name string, value uint, usage string) *uint
	FlagFloat64(name string, value float64, usage string) *float64

	GetBool(string) bool
	GetString(string) string
	GetDuration(string) time.Duration
	GetInt(string) int
	GetUint(string) uint
	GetFloat64(string) float64
}
