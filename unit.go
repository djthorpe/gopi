/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	UnitType uint
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type UnitConfig struct {
	Name string   // Unique name of the unit
	Type UnitType // Unit type
}

////////////////////////////////////////////////////////////////////////////////
// CONSTS

const (
	UNIT_NONE UnitType = iota
	UNIT_LOGGER
	UNIT_TIMER
	UNIT_MAX = UNIT_TIMER
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func UnitRegister(unit UnitConfig) {
	fmt.Printf("UnitRegister=%v\n", unit)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v UnitType) String() string {
	switch v {
	case UNIT_NONE:
		return "UNIT_NONE"
	case UNIT_LOGGER:
		return "UNIT_LOGGER"
	case UNIT_TIMER:
		return "UNIT_TIMER"
	default:
		return "[?? Invalid UnitType value]"
	}
}
