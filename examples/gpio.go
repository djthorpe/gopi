/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"fmt"
	"os"
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

	// Open up the GPIO interface
	gpio, err := rpi.NewGPIO()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	defer gpio.Close()
}
