/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"strconv"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	UnitType   uint
	ConfigFunc func(App) error
	NewFunc    func(App) (Unit, error)
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type UnitConfig struct {
	Name     string   // Unique name of the unit
	Type     UnitType // Unit type
	Requires []string // Unit dependencies
	Config   ConfigFunc
	New      NewFunc

	edges *unitArray
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

////////////////////////////////////////////////////////////////////////////////
// PRIVATE GLOBAL VARIABLES

var (
	unitMutex   sync.Mutex
	unitNameMap = map[string]UnitType{
		"logger": UNIT_LOGGER, // Logging
		"timer":  UNIT_TIMER,  // Timer
		"bus":    UNIT_BUS,    // Event Bus
	}
	unitByName map[string]*UnitConfig
	unitByType map[UnitType]*UnitConfig
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func UnitReset() {
	unitMutex.Lock()
	defer unitMutex.Unlock()
	unitByName = make(map[string]*UnitConfig)
	unitByType = make(map[UnitType]*UnitConfig)
}

func UnitRegister(unit UnitConfig) {
	// Reset
	if unitByName == nil || unitByType == nil {
		UnitReset()
	}

	// Lock
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
		if _, exists := unitByName[unit.Name]; exists {
			panic(fmt.Errorf("Duplicate Unit name registered: %v", unit))
		} else {
			unitByName[unit.Name] = &unit
		}
	}
	// Register by type
	if unit.Type != UNIT_NONE {
		if _, exists := unitByType[unit.Type]; exists {
			panic(fmt.Errorf("Duplicate Unit type registered: %v", unit))
		} else {
			unitByType[unit.Type] = &unit
		}
	}
}

func UnitsByType(unitType UnitType) []*UnitConfig {
	if unit, ok := unitByType[unitType]; ok {
		return []*UnitConfig{unit}
	} else {
		return nil
	}
}

func UnitsByName(unitName string) []*UnitConfig {
	if unitName == "" {
		units := make([]*UnitConfig, 0, len(unitByName))
		for _, unit := range unitByName {
			units = append(units, unit)
		}
		return units
	} else if unitType, exists := unitNameMap[unitName]; exists {
		return UnitsByType(unitType)
	} else if unit, ok := unitByName[unitName]; ok {
		return []*UnitConfig{unit}
	} else {
		return nil
	}
}

// UnitWithDependencies returns units with all dependencies
// satisfied in reverse order of when they need to be configured
// and initialized
func UnitWithDependencies(unitNames ...string) ([]*UnitConfig, error) {
	// Check incoming parameters
	if len(unitNames) == 0 {
		return nil, ErrBadParameter.WithPrefix("unitName")
	}

	unresolved := &unitArray{}
	resolved := &unitArray{}

	// Iterate through the units adding the edges to each module
	for _, name := range unitNames {
		if name == "" {
			return nil, ErrBadParameter.WithPrefix("unitName")
		} else if units := UnitsByName(name); len(units) == 0 {
			return nil, ErrBadParameter.WithPrefix(fmt.Sprintf("Unit not registered with name: %v", strconv.Quote(name)))
		} else {
			for _, unit := range units {
				if err := resolveUnitDependencies(unit, unresolved, resolved); err != nil {
					return nil, err
				}
			}
		}
	}
	// Success
	return resolved.arr, nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func resolveUnitDependencies(unit *UnitConfig, unresolved, resolved *unitArray) error {
	// Resolve edges if this unit hasn't been seen yet
	if unit.edges == nil {
		if edges, err := addUnitEdges(unit); err != nil {
			return err
		} else {
			unit.edges = edges
		}
	}

	// Mark as unresolved
	unresolved.Append(unit)

	// Now resolve each edge as necesary
	for _, edge := range unit.edges.arr {
		if resolved.Contains(edge) == false {
			if unresolved.Contains(edge) {
				return fmt.Errorf("Circular module reference detected: %v => %v", unit.Name, edge.Name)
			}
			if err := resolveUnitDependencies(edge, unresolved, resolved); err != nil {
				return err
			}
		}
	}

	// Module has been seen and can be removed from unresolved
	resolved.Append(unit)
	unresolved.Remove(unit)

	// Return success
	return nil
}

func addUnitEdges(unit *UnitConfig) (*unitArray, error) {
	edges := &unitArray{}
	for _, name := range unit.Requires {
		if name == "" {
			return nil, ErrBadParameter.WithPrefix(fmt.Sprintf("%s: Invalid Requires", unit.Name))
		} else if depends := UnitsByName(name); len(depends) == 0 {
			return nil, ErrNotFound.WithPrefix(fmt.Sprintf("%s: Unit with name: %s", unit.Name, name))
		} else {
			edges.Append(depends...)
		}
	}
	return edges, nil
}

type unitArray struct {
	arr []*UnitConfig
}

func (this *unitArray) Append(units ...*UnitConfig) {
	if this.arr == nil {
		this.arr = make([]*UnitConfig, 0, len(units))
	}
	this.arr = append(this.arr, units...)
}

func (this *unitArray) Remove(unit *UnitConfig) {
	for i, u := range this.arr {
		if u == unit {
			this.arr = append(this.arr[:i], this.arr[i+1:]...)
			break
		}
	}
}

func (this *unitArray) Contains(unit *UnitConfig) bool {
	for _, u := range this.arr {
		if u == unit {
			return true
		}
	}
	return false
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
