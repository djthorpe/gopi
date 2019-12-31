/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi_test

import (
	"fmt"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

func Test_Unit_000(t *testing.T) {
	t.Log("Test_Unit_000")
}

func Test_Unit_001(t *testing.T) {
	gopi.UnitReset()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	gopi.UnitRegister(gopi.UnitConfig{})
	t.Error("Expected panic")
}

func Test_Unit_002(t *testing.T) {
	gopi.UnitReset()
	gopi.UnitRegister(gopi.UnitConfig{
		Type: gopi.UNIT_LOGGER,
	})
}

func Test_Unit_003(t *testing.T) {
	gopi.UnitReset()
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "test",
	})
}

func Test_Unit_004(t *testing.T) {
	gopi.UnitReset()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "logger",
	})
	t.Error("Expected panic")
}

func Test_Unit_005(t *testing.T) {
	gopi.UnitReset()

	// Two modules cannot have the same name
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered:", r)
		}
	}()
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "test",
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "test",
		Type: gopi.UNIT_LOGGER,
	})
	t.Error("Expected panic")
}

func Test_Unit_006(t *testing.T) {
	gopi.UnitReset()
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "test_bus",
		Type: gopi.UNIT_BUS,
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "test_logger",
		Type: gopi.UNIT_LOGGER,
	})

	if units := gopi.UnitsByType(gopi.UNIT_BUS); len(units) != 1 {
		t.Error("Unexpected return value from UnitsByType")
	} else if units[0].Name != "test_bus" {
		t.Error("Unexpected return value from UnitsByType")
	}

	if units := gopi.UnitsByType(gopi.UNIT_LOGGER); len(units) != 1 {
		t.Error("Unexpected return value from UnitsByType")
	} else if units[0].Name != "test_logger" {
		t.Error("Unexpected return value from UnitsByType")
	}

	if units := gopi.UnitsByName("test_logger"); len(units) != 1 {
		t.Error("Unexpected return value from UnitsByName")
	} else if units[0].Type != gopi.UNIT_LOGGER {
		t.Error("Unexpected return value from UnitsByName")
	}

	if units := gopi.UnitsByName("test_bus"); len(units) != 1 {
		t.Error("Unexpected return value from UnitsByName")
	} else if units[0].Type != gopi.UNIT_BUS {
		t.Error("Unexpected return value from UnitsByName")
	}

}

func Test_Unit_007(t *testing.T) {
	gopi.UnitReset()
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "A",
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "B",
		Requires: []string{"A"},
	})

	if units, err := gopi.UnitWithDependencies("A"); err != nil {
		t.Error(err)
	} else if len(units) != 1 {
		t.Error("Expected one unit, got", units)
	} else if units[0].Name != "A" {
		t.Error("Expected unit A, got", units[0])
	}

	if units, err := gopi.UnitWithDependencies("B"); err != nil {
		t.Error(err)
	} else if len(units) != 2 {
		t.Error("Expected two units, got", units)
	} else if units[0].Name != "A" {
		t.Error("Expected unit A, got", units[0])
	} else if units[1].Name != "B" {
		t.Error("Expected unit B, got", units[1])
	}
}

func Test_Unit_008(t *testing.T) {
	gopi.UnitReset()
	gopi.UnitRegister(gopi.UnitConfig{
		Name: "A",
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "B",
		Requires: []string{"A"},
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "C",
		Requires: []string{"A", "B"},
	})

	if units, err := gopi.UnitWithDependencies("C"); err != nil {
		t.Error(err)
	} else if len(units) != 3 {
		t.Error("Expected three units, got", units)
	} else if units[2].Name != "C" {
		t.Error("Expected unit C, got", units[2])
	} else if units[1].Name != "B" {
		t.Error("Expected unit B, got", units[1])
	} else if units[0].Name != "A" {
		t.Error("Expected unit A, got", units[0])
	}
}

func Test_Unit_009(t *testing.T) {
	gopi.UnitReset()
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "A",
		Requires: []string{"B"},
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "B",
		Requires: []string{"A"},
	})
	gopi.UnitRegister(gopi.UnitConfig{
		Name:     "C",
		Requires: []string{"A", "B"},
	})

	if _, err := gopi.UnitWithDependencies("C"); err == nil {
		t.Error("Expected circular dependency error")
	} else {
		t.Log(err)
	}
}
