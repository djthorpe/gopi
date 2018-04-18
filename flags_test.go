package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
)

func TestFlags_000(t *testing.T) {
	// Create a configuration with debug
	flagset := gopi.NewFlags("test")
	if flagset == nil {
		t.Fatalf("Expected flagset != nil")
	}
	if flagset.Name() != "test" {
		t.Fatalf("Unexpected Name() return:", flagset.Name())
	}
}

func TestFlags_001(t *testing.T) {
	// Create a configuration with debug
	flagset := gopi.NewFlags("test")
	if flagset.Parsed() != false {
		t.Error("Unexpected Parsed() return value")
	}
	if err := flagset.Parse([]string{}); err != nil {
		t.Error("Unexpected Parse() error:", err)
	}
	if flagset.Parsed() != true {
		t.Error("Unexpected Parsed() return value")
	}
}

func TestFlags_002(t *testing.T) {
	// Create a configuration with debug
	flagset := gopi.NewFlags("test")
	if err := flagset.Parse([]string{"-help"}); err != gopi.ErrHelp {
		t.Error("Unexpected Parse() error:", err)
	}
}

func TestFlags_003(t *testing.T) {
	// Create a configuration with debug
	flagset := gopi.NewFlags("test")
	if err := flagset.Parse([]string{"a", "b"}); err != nil {
		t.Error("Unexpected Parse() error:", err)
	}
	args := flagset.Args()
	if len(args) != 2 || args[0] != "a" || args[1] != "b" {
		t.Error("Unexpected Args() return:", args)
	}
}

func TestFlags_004(t *testing.T) {
	// Create a configuration with debug
	flagset := gopi.NewFlags("test")
	flagset.FlagBool("test", false, "test argument")
	if err := flagset.Parse([]string{"-test"}); err != nil {
		t.Error("Unexpected Parse() error:", err)
	}
	value, exists := flagset.GetBool("test")
	if value == false || exists == false {
		t.Error("Unexpected GetBool() return")
	}
}

func TestFlags_005(t *testing.T) {
	// Create a configuration with debug
	flagset := gopi.NewFlags("test")
	flagset.FlagBool("test", true, "test argument")
	if err := flagset.Parse([]string{""}); err != nil {
		t.Error("Unexpected Parse() error:", err)
	}
	value, exists := flagset.GetBool("test")
	if value == false || exists == true {
		t.Error("Unexpected GetBool() return")
	}
}
