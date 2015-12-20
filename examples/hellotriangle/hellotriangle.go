package main

import (
	"github.com/djthorpe/gopi/rpi/egl"
	"log"
)


func main() {
	egl.BCMHostInit()

	// Initalize display
	display := egl.GetDisplay()
	if ok,err := egl.Initialize(display,nil,nil); err != nil {
		log.Fatalf("Unable to initalize display: %v",err)
	}

	// Terminate display
	if ok,err := egl.Terminate(display); err != nil {
		log.Fatalf("Unable to terminate display: %v",err)
	}
}


