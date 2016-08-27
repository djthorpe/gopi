package main

import (
	"flag"
	"fmt"
	"github.com/djthorpe/gopi/rpi"
	"github.com/djthorpe/gopi/rpi/displaymanx"
	"github.com/djthorpe/gopi/rpi/egl"
	"log"
	"runtime"
)

var (
	displayFlag = flag.Uint("display", 0, "Display Number")
)

func GetConfigAttributeValue(display egl.Display, config egl.Config, name int32) int32 {
	value, err := egl.GetConfigAttribute(display, config, name)
	if err != nil {
		log.Fatalf("GetConfigAttribute: %v", err)
	}
	return value
}

func PrintConfig(display egl.Display, config egl.Config) {
	fmt.Println("Configuration ID ", GetConfigAttributeValue(display, config, egl.EGL_CONFIG_ID))
	fmt.Println("\t  EGL_BUFFER_SIZE: ", GetConfigAttributeValue(display, config, egl.EGL_BUFFER_SIZE))
	fmt.Println("\t     EGL_RED_SIZE: ", GetConfigAttributeValue(display, config, egl.EGL_RED_SIZE))
	fmt.Println("\t   EGL_GREEN_SIZE: ", GetConfigAttributeValue(display, config, egl.EGL_GREEN_SIZE))
	fmt.Println("\t    EGL_BLUE_SIZE: ", GetConfigAttributeValue(display, config, egl.EGL_BLUE_SIZE))
	fmt.Println("\t   EGL_ALPHA_SIZE: ", GetConfigAttributeValue(display, config, egl.EGL_ALPHA_SIZE))
	fmt.Println("\t   EGL_DEPTH_SIZE: ", GetConfigAttributeValue(display, config, egl.EGL_DEPTH_SIZE))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rpi.BCMHostInit()
	flag.Parse()

	// Initalize display
	display := egl.GetDisplay()
	if err := egl.Initialize(display, nil, nil); err != nil {
		log.Fatalf("Initialize: %v", err)
	}

	// Choose a configuration which matches
	attr := []int32{
		egl.EGL_RED_SIZE, 8,
		egl.EGL_GREEN_SIZE, 8,
		egl.EGL_BLUE_SIZE, 8,
		egl.EGL_ALPHA_SIZE, 8,
		egl.EGL_SURFACE_TYPE, egl.EGL_WINDOW_BIT,
		egl.EGL_NONE,
	}

	configs, err := egl.ChooseConfig(display, attr)
	if err != nil {
		log.Fatalf("ChooseConfig: %v", err)
	}
	if len(configs) == 0 {
		log.Fatalf("ChooseConfig: Failed to choose appropriate graphics mode")
	}

	// Print all Configs
	for i := 0; i < len(configs); i++ {
		PrintConfig(display, configs[i])
	}

	// Bind API
	if err := egl.BindAPI(egl.EGL_OPENGL_ES_API); err != nil {
		log.Fatalf("BindAPI: %v", err)
	}

	// Context
	ctxAttr := []int32{
		egl.EGL_CONTEXT_CLIENT_VERSION, 2,
		egl.EGL_NONE,
	}
	context, err := egl.CreateContext(display, configs[0], egl.EGL_NO_CONTEXT, &ctxAttr[0])
	if err != nil {
		log.Fatalf("CreateContext: %v", err)
	}

	screen_width, screen_height := rpi.GraphicsGetDisplaySize(uint16(*displayFlag))
	fmt.Println("Screen size = (", screen_width, ",", screen_height, ")")

	srcRect := displaymanx.Rect{0, 0, screen_width, screen_height}
	dstRect := displaymanx.Rect{0, 0, screen_width << 16, screen_height << 16}

	// Bind to actual screen
	dispmanx_display := displaymanx.DisplayOpen(uint32(*displayFlag))
	dispmanx_update := displaymanx.UpdateStart(0) /* priority */
	dispmanx_element := displaymanx.ElementAdd(dispmanx_update, dispmanx_display, 0, &dstRect, 0, &srcRect, displaymanx.DISPMANX_PROTECTION_NONE, nil, nil, 0)
	window := displaymanx.Window{dispmanx_element, screen_width, screen_height}

	// do something here
	displaymanx.UpdateSubmitSync(dispmanx_update)

	// make surface
	surface, err := egl.CreateWindowSurface(display, configs[0], window, nil)
	if err != nil {
		log.Fatalf("CreateWindowSurface: %v", err)
	}

	// connect the context to the surface
	if err := egl.MakeCurrent(display, surface, surface, context); err != nil {
		log.Fatalf("MakeCurrent: %v", err)
	}

	// TODO

	// Destroy surface
	if err := egl.DestroySurface(display, surface); err != nil {
		log.Fatalf("DestroySurface: %v", err)
	}

	// Destroy context
	if err := egl.DestroyContext(display, context); err != nil {
		log.Fatalf("DestroyContext: %v", err)
	}

	// Terminate display
	if err := egl.Terminate(display); err != nil {
		log.Fatalf("Terminate: %v", err)
	}
}
