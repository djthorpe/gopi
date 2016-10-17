/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

import (
	rpi "../device/rpi"
)

////////////////////////////////////////////////////////////////////////////////

var (
	flagBus = flag.Uint("bus", 1, "I2C Bus")
)

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Open up the RaspberryPi interface
	rpi, err := rpi.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(-1)
	}
	defer rpi.Close()

	// Set flag usage, parse flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", path.Base(os.Args[0]))
		flag.PrintDefaults()
		return
	}
	flag.Parse()

	// Open up the I2C interface
	i2c, err := rpi.NewI2C(*flagBus)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer i2c.Close()
}
