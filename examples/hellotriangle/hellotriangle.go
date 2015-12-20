package main

import (
	"github.com/djthorpe/gopi/rpi/egl"
	"log"
	"os"
)


func main() {
	egl.BCMHostInit()

	// Initalize display
	display := egl.GetDisplay(egl.DEFAULT_DISPLAY)
	if ok := egl.Initialize(display,nil,nil); ok != true {
		log.Errorf("Unable to initalize display")
		return -1
	}

	// Terminate display
	if ok := egl.Terminate(display); ok != true {
		log.Errorf("Unable to terminate display")
		return -1
	}

	return 0
}


