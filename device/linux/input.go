/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Empty input configuration
type Input struct{}

// Driver of multiple input devices
type InputDriver struct {
	log     *util.LoggerDevice // logger
	devices []hw.InputDevice   // input devices
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new Input object, returns error if not possible
func (config Input) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<linux.Input>Open")

	// create new GPIO driver
	this := new(InputDriver)

	// Set logging & devices
	this.log = log
	this.devices = make([]hw.InputDevice,0)

	// TODO

	// success
	return this, nil
}

// Close Input driver
func (this *InputDriver) Close() error {
	this.log.Debug("<linux.Input>Close")

	for _, device := range this.devices {
		if err := device.Close(); err != nil {
			return err
		}
	}

	return nil
}
