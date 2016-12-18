/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Empty input configuration
type Input struct { }

type InputDriver struct {
	log   *util.LoggerDevice // logger
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new Input object, returns error if not possible
func (config Input) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<linux.Input>Open")

	var err error

	// create new GPIO driver
	this := new(InputDriver)

	// Set logging & device
	this.log = log

	// success
	return this, nil
}

// Close Input driver
func (this *InputDriver) Close() error {
	this.log.Debug("<linux.Input>Close")

	return nil
}

// Strinfigy InputDriver object
func (this *InputDriver) String() string {
	return fmt.Sprintf("<linux.Input>{ }")
}

