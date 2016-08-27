/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

import (
	"os"
)

////////////////////////////////////////////////////////////////////////////////

type I2C struct {
	file *os.File
}

////////////////////////////////////////////////////////////////////////////////

// Create new I2C object, returns error if not possible
func (rpi *RaspberryPi) NewI2C(bus uint,addr uint8) (*I2C, error) {
	file, err := os.OpenFile(dev,os.O_RDWR,0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(file.Fd(),i2c_SLAVE,uintptr(addr)); err != nil {
		return nil, err
	}
	return &I2C{ file }, nil
}

func (this *I2C) Close() {
	this.file.Close()
}

