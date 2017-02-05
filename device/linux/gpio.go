/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"fmt"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GPIO struct { }

type gpioDriver struct {
	log      *util.LoggerDevice // logger
	pins     map[hw.GPIOPin]uint // map of logical to physical pins
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GPIO_DEV        = "/sys/class/gpio"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config GPIO) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<linux.GPIO>Open")

	// create new GPIO driver
	this := new(gpioDriver)

	// Set logging & device
	this.log = log

	// TODO

	// success
	return this, nil
}

// Close GPIO connection
func (this *gpioDriver) Close() error {
	this.log.Debug("<linux.GPIO>Close")

	return nil
}

// Stringfigy GPIO
func (this *gpioDriver) String() string {
	return fmt.Sprintf("<linux.GOPI>{ }")
}

