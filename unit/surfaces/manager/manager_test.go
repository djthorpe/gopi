// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package manager_test

import (
	"errors"
	"fmt"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	manager "github.com/djthorpe/gopi/v2/unit/surfaces/manager"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN DISPLAY

var (
	Display rpi.DXDisplayHandle
)

func init() {
	rpi.DXInit()
	if display, err := rpi.DXDisplayOpen(rpi.DXDisplayId(0)); err != nil {
		panic(fmt.Sprint(err))
	} else {
		Display = display
	}
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Manager_000(t *testing.T) {
	t.Log("Test_Manager_000")
}

func Test_Manager_001(t *testing.T) {
	// Display parameter is required
	if _, err := gopi.New(manager.Config{}, NewLogger(t)); errors.Is(err, gopi.ErrBadParameter) == false {
		t.Error("Unexpected error return", err)
	}
	// Should work
	if manager, err := gopi.New(manager.Config{Display}, NewLogger(t)); err != nil {
		t.Error(err)
	} else if err := manager.Close(); err != nil {
		t.Error(err)
	} else {
		t.Log(manager)
	}
}

func Test_Manager_002(t *testing.T) {
	mm, err := gopi.New(manager.Config{Display}, NewLogger(t))
	if err != nil {
		t.Fatal(err)
	}

	if bm1, err := mm.(manager.Manager).NewBitmap(gopi.Size{1, 1}, 0); err != nil {
		t.Error(err)
	} else if bm2, err := mm.(manager.Manager).NewBitmap(gopi.Size{2, 2}, 0); err != nil {
		t.Error(err)
	} else if bm3, err := mm.(manager.Manager).NewBitmap(gopi.Size{3, 3}, 0); err != nil {
		t.Error(err)
	} else if err := mm.(manager.Manager).ReleaseBitmap(bm1); err != nil {
		t.Error(err)
	} else if err := mm.(manager.Manager).ReleaseBitmap(bm2); err != nil {
		t.Error(err)
	} else if err := mm.(manager.Manager).ReleaseBitmap(bm3); err != nil {
		t.Error(err)
	}

	if err := mm.Close(); err != nil {
		t.Fatal(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// LOGGER

type logger struct{ testing.T }

func NewLogger(t *testing.T) gopi.Logger {
	return &logger{*t}
}

func (this *logger) Clone(string) gopi.Logger {
	return this
}
func (this *logger) Name() string {
	return this.T.Name()
}
func (this *logger) Error(err error) error {
	this.T.Error(err)
	return err
}
func (this *logger) Warn(args ...interface{}) {
	this.T.Log(args...)
}
func (this *logger) Info(args ...interface{}) {
	this.T.Log(args...)
}
func (this *logger) Debug(args ...interface{}) {
	this.T.Log(args...)
}
func (this *logger) IsDebug() bool {
	return true
}
func (this *logger) Close() error {
	return nil
}
func (this *logger) String() string {
	return this.Name()
}
