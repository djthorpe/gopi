/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

import (
	"fmt"
	"os"
)

////////////////////////////////////////////////////////////////////////////////

type I2C struct {
	file *os.File
}

////////////////////////////////////////////////////////////////////////////////

const (
	I2C_DEV = "/dev/i2c"
)

////////////////////////////////////////////////////////////////////////////////

// Create new I2C object, returns error if not possible
func (rpi *RaspberryPi) NewI2C(bus uint) (*I2C, error) {
	device, err := getDevice(bus)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(device, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	/*	if err := ioctl(file.Fd(),i2c_SLAVE,uintptr(addr)); err != nil {
		return nil, err
	}*/
	return &I2C{file}, nil
}

func (this *I2C) Close() {
	this.file.Close()
}

////////////////////////////////////////////////////////////////////////////////

func getDevice(bus uint) (string, error) {
	device := fmt.Sprintf("%s-%v", I2C_DEV, bus)
	if err := isReadablePath(device); err != nil {
		return "", err
	}
	return device, nil
}
