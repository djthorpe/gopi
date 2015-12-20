package main

import (
	"github.com/djthorpe/gopi/rpi/egl"
	"log"
)


func main() {
	egl.BCMHostInit()

	// Initalize display
	display := egl.GetDisplay()
	if ok := egl.Initialize(display,nil,nil); ok != true {
		log.Fatalf("Unable to initalize display")
	}

	// Terminate display
	if ok := egl.Terminate(display); ok != true {
		log.Fatalf("Unable to terminate display")
	}

	return 0
}


