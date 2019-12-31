/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MainCommandFunc func(App, []string) error
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type App interface {
	Run() int // Run application

	Log() Logger  // Return logger unit
	Timer() Timer // Return timer unit
	Bus() Bus     // Return event bus unit

	Unit(string) Unit    // Return singular unit for name
	Units(string) []Unit // Return multiple units for name
}
