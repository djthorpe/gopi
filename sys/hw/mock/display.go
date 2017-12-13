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

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Display struct {
	Display uint
}

type display struct {
	id  uint
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Display) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.mock.Display.Open{ id=%v }", config.Display)

	this := new(display)
	this.log = logger
	this.id = config.Display

	// Success
	return this, nil
}

// Close
func (this *display) Close() error {
	this.log.Debug("sys.mock.Display.Close{ id=%v }", this.id)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return display number
func (this *display) Display() uint {
	return 0
}

// Return size
func (this *display) Size() (uint32, uint32) {
	return 0, 0
}

// Return pixels-per-inch
func (this *display) PixelsPerInch() uint32 {
	return 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *display) String() string {
	return fmt.Sprintf("sys.mock.Display{ id=%v }", this.id)
}
