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
}