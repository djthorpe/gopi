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
	"os"
	"sync"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SPI struct {
	Bus    uint
	Channel uint
}

type spiDriver struct {
	log   *util.LoggerDevice // logger
	dev   *os.File
	bus   uint
	channel uint
	lock  sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SPI_DEV                   = "/dev/spidev"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new SPI object, returns error if not possible
func (config SPI) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<linux.SPI>Open")

	var err error

	// create new GPIO driver
	this := new(spiDriver)
	this.bus = config.Bus
	this.channel = config.Channel

	// Set logging & device
	this.log = log

	// Open the device
	this.dev, err = os.OpenFile(fmt.Sprintf("%v%v.%v",SPI_DEV,this.bus,this.channel),os.O_RDWR,0)
	if err != nil {
		return nil, err
	}

	// success
	return this, nil
}

// Close SPI connection
func (this *spiDriver) Close() error {
	this.log.Debug("<linux.SPI>Close")

	err := this.dev.Close()
	this.dev = nil
	return err
}

// Strinfigy SPI driver
func (this *spiDriver) String() string {
	return fmt.Sprintf("<linux.SPI>{ bus=%v channel=%v }", this.bus, this.channel)
}

