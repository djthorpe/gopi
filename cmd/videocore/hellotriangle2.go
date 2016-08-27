package main

import (
	"flag"
	"github.com/djthorpe/gopi/rpi"
	"github.com/djthorpe/gopi/rpi/egl"
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

	// Initalize display
	attr := []int32{
		egl.EGL_RED_SIZE, 8,
		egl.EGL_GREEN_SIZE, 8,
		egl.EGL_BLUE_SIZE, 8,
		egl.EGL_ALPHA_SIZE, 8,
		egl.EGL_SURFACE_TYPE, egl.EGL_WINDOW_BIT,
		egl.EGL_NONE,
	}
	graphics, err := egl.New(uint16(*displayFlag), attr)
	if err != nil {
		log.Fatalf("New: %v", err)
	}

	// TODO

	// Terminate display
	if err := graphics.Terminate(); err != nil {
		log.Fatalf("Terminate: %v", err)
	}
}
