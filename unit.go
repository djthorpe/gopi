/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	UnitType   uint
	ConfigFunc func(App) error
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type UnitConfig struct {
	Name     string   // Unique name of the unit
	Type     UnitType // Unit type
	Requires []string // Unit dependencies
	Config   ConfigFunc
}

////////////////////////////////////////////////////////////////////////////////
// CONSTS

const (
	UNIT_NONE UnitType = iota
	UNIT_LOGGER
	UNIT_TIMER
	UNIT_BUS
	UNIT_MAX = UNIT_BUS
)

var (
	unitMutex   sync.Mutex
	unitNameMap = map[string]UnitType{
		"logger": UNIT_LOGGER, // Logging
		"timer":  UNIT_TIMER,  // Timer
		"bus":    UNIT_BUS,    // Event Bus
	}
	unitByName = make(map[string]UnitConfig)
	unitByType = make(map[UnitType]UnitConfig)
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func UnitRegister(unit UnitConfig) {
	unitMutex.Lock()
	defer unitMutex.Unlock()

	// Satisfy unit.Type or unit.Name
	if unit.Type == UNIT_NONE {
		if unit.Name == "" {
			panic(fmt.Errorf("Unit name cannot be empty when type is NONE: %v", unit))
		}
	}
	// Unit cannot be a reserved name
	if _, exists := unitNameMap[unit.Name]; exists {
		panic(fmt.Errorf("Unit name uses reserved word: %v", unit))
	}
	// Register by name
	if unit.Name != "" {
		if _, exists := unitNameMap[unit.Name]; exists {
			panic(fmt.Errorf("Duplicate Unit registered: %v", unit))
		} else {
			unitByName[unit.Name] = unit
		}
	}
	// Register by type
	if unit.Type != UNIT_NONE {
		if _, exists := unitByType[unit.Type]; exists {
			panic(fmt.Errorf("Duplicate Unit registered: %v", unit))
		} else {
			unitByType[unit.Type] = unit
		}
	}
}

func UnitsByType(unitType UnitType) []UnitConfig {
	unitMutex.Lock()
	defer unitMutex.Unlock()

	if unit, ok := unitByType[unitType]; ok {
		return []UnitConfig{unit}
	} else {
		return nil
	}
}

func UnitsByName(unitName string) []UnitConfig {
	unitMutex.Lock()
	defer unitMutex.Unlock()

	if unit, ok := unitByName[unitName]; ok {
		return []UnitConfig{unit}
	} else {
		return nil
	}
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
	case UNIT_BUS:
		return "UNIT_BUS"
	default:
		return "[?? Invalid UnitType value]"
	}
}

func (u UnitConfig) String() string {
	return fmt.Sprintf("<gopi.Unit name=%s type=%s requires=%s>", u.Name, u.Type, u.Requires)
}
