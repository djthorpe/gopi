/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"flag"
)

import (
	"../" /* import "github.com/djthorpe/gopi" */
	"../util" /* import "github.com/djthorpe/gopi/util" */
	"../device/rpi" /* import "github.com/djthorpe/gopi/device/rpi" */
)

var (
	flagDisplay = flag.Uint("display", 0,"Display number")
)

func main() {
	// Parse flags
	flag.Parse()
	
	// Create the logger
	log, err := util.Logger(util.StderrLogger{ })
	if err != nil {
		panic("Can't open logging interface")
	}
	defer log.Close()

	// Set logging level
	log.SetLevel(util.LOG_ANY)

	// Open the Raspberry Pi device
	device, err := gopi.Open(rpi.Device{ },log)
	if err != nil {
		log.Fatal("%v",err)
		return
	}
	defer device.Close()

	// Open the display
	display, err := device.Display(rpi.DXDisplayConfig{ uint16(*flagDisplay) })
	if err != nil {
		log.Fatal("%v",err)
		return
	}
	defer display.Close()

	// Create a resource
	resource, err := display.CreateResource(display.GetSize())
	if err != nil {
		log.Fatal("%v",err)
		return
	}
	defer display.CloseResource(resource)

	log.Info("Device=%v",device)
	log.Info("Display=%v",display)
	log.Info("Resource=%v",resource)
}
