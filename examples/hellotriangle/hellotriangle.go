package main

import (
	"github.com/djthorpe/gopi/rpi/egl"
	"log"
)


func main() {
	log.Println("BCM Host Init")
	egl.BCMHostInit()
	log.Println("Done")
}


