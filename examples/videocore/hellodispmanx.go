package main

import (
	"flag"
	"github.com/djthorpe/gopi/rpi"
	"github.com/djthorpe/gopi/rpi/dispmanx"
	"log"
	"runtime"
)

var (
	displayFlag = flag.Uint("display", 0, "Display Number")
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rpi.BCMHostInit()
	flag.Parse()

	display, err := dispmanx.New(uint32(*displayFlag))
	if err != nil {
		log.Fatalf("New: %v", err)
	}
	defer display.Terminate()

	// Start update
	if err := display.StartUpdate(10); err != nil {
		log.Fatalf("Start Update: %v", err)
	}

	// End update
	if err := display.EndUpdate(); err != nil {
		log.Fatalf("End Update: %v", err)
	}

}
