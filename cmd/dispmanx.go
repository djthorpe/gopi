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
	"time"
)

import (
	rpi "../device/rpi"
)

var (
	flagDisplay = flag.Uint("display", 0, "Display Number")
)

func main() {
	// Flags
	flag.Parse()

	// Open up the RaspberryPi interface
	pi, err := rpi.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	defer pi.Close()

	// VideoCore
	vc, err := pi.NewVideoCore(uint16(*flagDisplay))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	defer vc.Close()

	fmt.Printf("Display=%v Size=%v\n", vc.GetDisplay(), vc.GetSize())

	info, err := vc.GetModeInfo()
	fmt.Printf("ModeInfo=%v\n", info)

	// Create a background resource
	background, err := vc.CreateResource(rpi.VC_IMAGE_RGB565, rpi.Size{1, 1})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}

	// Start an update
	handle, err := vc.UpdateBegin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	vc.SetBackgroundColor(handle, rpi.Color{255, 0, 0})
	vc.UpdateSubmit(handle)

	// Delete the background
	vc.DeleteResource(background)

	// Wait...
	time.Sleep(time.Second * 3)
}
