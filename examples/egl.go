/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example demonstrates the use of OpenVG for vector graphics
package main

import (
	"flag"
	"os"
	"fmt"
)

import (
	rpi "../device/rpi"
	khronos "../khronos"
	util "../util"
)

////////////////////////////////////////////////////////////////////////////////

var (
	flagDisplay = flag.String("display", "lcd", "Display")
)
////////////////////////////////////////////////////////////////////////////////

func main() {
	// Flags
	flag.Parse()

	// Create logger
	logger := new(util.StderrLogger)

	// Retrieve display
	display, ok := rpi.Displays[*flagDisplay]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: Invalid display name")
		return
	}

	// Open up the RaspberryPi interface
	pi, err := rpi.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer pi.Close()

	// VideoCore
	vc, err := pi.NewVideoCore(display)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		return
	}
	defer vc.Close()

	// EGL
	egl, err := khronos.Open(&rpi.EGL{ vc, logger })
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer egl.Close()

	// Create a window
	window, err := egl.CreateWindow("OpenVG",&khronos.EGLSize{ 100, 100 },&khronos.EGLPoint{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// DO SOMETHING HERE
	fmt.Println(window)
/*

	if err := egl.Do(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	*/

}
