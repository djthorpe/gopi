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
		log.Printf("Unable to initalize display")
		return -1
	}

	// Terminate display
	if ok := egl.Terminate(display); ok != true {
		log.Printf("Unable to terminate display")
		return -1
	}

	return 0
}


