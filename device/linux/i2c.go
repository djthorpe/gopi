/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"os"
	"fmt"
	"syscall"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2C struct {
	Device gopi.HardwareDriver
	Bus    uint
	Slave  uint8
}

type I2CDriver struct {
	log      *util.LoggerDevice // logger
	bus      uint
	slave    uint8
	dev      *os.File
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	I2C_DEV = "/dev/i2c"
)

const (
	// i2c ioctl commands
	I2C_RETRIES      = 0x0701  /* number of times a device address should be polled when not acknowledging */
	I2C_TIMEOUT      = 0x0702  /* set timeout in units of 10 ms */
	I2C_SLAVE        = 0x0703  /* Use this slave address */
	I2C_SLAVE_FORCE  = 0x0706  /* Use this slave address, even if it is already in use by a driver! */
	I2C_TENBIT       = 0x0704  /* 0 for 7 bit addrs, != 0 for 10 bit */
	I2C_FUNCS        = 0x0705  /* Get the adapter functionality mask */
	I2C_RDWR         = 0x0707  /* Combined R/W transfer (one STOP only) */
	I2C_PEC          = 0x0708  /* != 0 to use PEC with SMBus */
	I2C_SMBUS        = 0x0720  /* SMBus transfer */
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new I2C object, returns error if not possible
func (config I2C) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<linux.I2C>Open")

	var err error

	// create new GPIO driver
	this := new(I2CDriver)

	// Set logging & device
	this.log = log
	this.bus = config.Bus
	this.slave = config.Slave

	// Open the /dev/mem and provide offset & size for accessing memory
	this.dev, err = i2cOpenDevice(config.Bus)
	if err != nil {
		return nil, err
	}

	// Set address
	if err := i2cIoctl(this.dev.Fd(), I2C_SLAVE, uintptr(this.slave)); err != nil {
		return nil, err
	}

	// success
	return this, nil
}

// Close I2C connection
func (this *I2CDriver) Close() error {
	this.log.Debug("<linux.I2C>Close")

	err := this.dev.Close()
	this.dev = nil
	return err
}

// Strinfigy I2C object
func (this *I2CDriver) String() string {
	return fmt.Sprintf("<linux.I2C>{ bus=%v slave=0x%02X }", this.bus, this.slave)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func i2cOpenDevice(bus uint) (*os.File,error) {
	var file *os.File
	var err error

	if file, err = os.OpenFile(fmt.Sprintf("%v-%v",I2C_DEV,bus),os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	}
	return file, nil
}

func i2cIoctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}





