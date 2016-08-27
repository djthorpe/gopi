/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"os"
	"fmt"
)

import (
	rpi "../device/rpi"
)

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Open up the RaspberryPi interface
	rpi, err := rpi.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	defer rpi.Close()

	// Open up the I2C interface
	i2c, err := rpi.NewI2C()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	defer i2c.Close()
}
