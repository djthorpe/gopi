/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package mock

import (
	"fmt"
	"strings"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Hardware struct{}

type hardware struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register hardware
	gopi.RegisterModule(gopi.Module{
		Name: "hardware/mock",
		Type: gopi.MODULE_TYPE_HARDWARE,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Hardware{}, app.Logger)
		},
	})
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Hardware) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.mock.Hardware.Open{  }")

	this := new(hardware)
	this.log = logger

	// Success
	return this, nil
}

// Close
func (this *hardware) Close() error {
	this.log.Debug("sys.mock.Hardware.Close{ }")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetName returns the name of the hardware (ie, mock, mac, linux, rpi, etc)
func (this *hardware) Name() string {
	return "mock/hw"
}

// SerialNumber returns the serial number of the hardware, if available
func (this *hardware) SerialNumber() string {
	return strings.ToUpper("<SERIAL_NUMBER>")
}

// Return the number of displays which can be opened
func (this *hardware) NumberOfDisplays() uint {
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *hardware) String() string {
	return fmt.Sprintf("sys.mock.Hardware{ name=%v serial=%v displays=%v }", this.Name(), this.SerialNumber(), this.NumberOfDisplays())
}
