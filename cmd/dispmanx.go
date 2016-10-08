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

	// Start an update
	handle, err := vc.UpdateBegin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}

	// Add an element onto the screen
	dest := Rectangle{ Point{ 0,0 }, info.Size }
	src := Rectangle{ Point{ 0,0 }, info.Size }
	element := vc.ElementAdd(handle,Layer(0),dest,0,src,rpi.DISPMANX_PROTECTION_NONE,Alpha(0),Clamp(0),Transform(0))

	vc.SetBackgroundColor(handle, rpi.Color{ 0, 200, 0})
	vc.UpdateSubmit(handle)

	// Wait...
	time.Sleep(time.Second * 3)
}
