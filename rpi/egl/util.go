/*
Go Language Raspberry Pi Interface
(c) Copyright David Thorpe 2016
All Rights Reserved

For Licensing and Usage information, please see LICENSE.md
*/
package egl

import (
	"github.com/djthorpe/gopi/rpi"
	"github.com/djthorpe/gopi/rpi/dispmanx"
)

////////////////////////////////////////////////////////////////////////////////

type EGL struct {
	major,minor int32 // EGL version
	displayn uint16 // Display number
	displayw,displayh uint32 // Display size
    displayp Display // Display handle
    contextp Context // Context handle
    surfacep Surface // Surface handle
}

////////////////////////////////////////////////////////////////////////////////

func New(display uint16,attr []int32) (*EGL,error) {

	// create an EGL variable
    this := new(EGL)
	this.displayn = display
	this.displayp = EGL_NO_DISPLAY
	this.contextp = EGL_NO_CONTEXT
	this.surfacep = EGL_NO_SURFACE

    // Initalize display
    this.displayp = GetDisplay()
    if err := Initialize(this.displayp,&this.major,&this.minor); err != nil {
        return nil,err
    }

	// Choose configuration
	configs, err := ChooseConfig(this.displayp, attr)
	if err != nil {
		return nil,err
	}
	if len(configs) == 0 {
		return nil,ErrorInvalidGraphicsConfiguration
	}

	// Bind API
	if err := BindAPI(EGL_OPENGL_ES_API); err != nil {
		return nil,err
	}

	// Context
	ctxAttr := []int32{
		EGL_CONTEXT_CLIENT_VERSION, 2,
		EGL_NONE,
	}	
	this.contextp,err = CreateContext(this.displayp,configs[0],EGL_NO_CONTEXT,&ctxAttr[0])
	if err != nil {
		return nil,err
	}

	// Get display width and height
	this.displayw,this.displayh = rpi.GraphicsGetDisplaySize(this.displayn)

	// Create window
	srcRect := dispmanx.Rect{ 0,0,this.displayw,this.displayh }
	dstRect := dispmanx.Rect{ 0,0,this.displayw << 16,this.displayh << 16 }

	// Bind to screen
	dispmanx_display := dispmanx.DisplayOpen(uint32(this.displayn))
	dispmanx_update := dispmanx.UpdateStart(0) /* priority */
	dispmanx_element := dispmanx.ElementAdd(dispmanx_update,dispmanx_display,0,&dstRect,0,&srcRect,dispmanx.DISPMANX_PROTECTION_NONE,nil,nil,0)
	window := dispmanx.Window{ dispmanx_element,this.displayw,this.displayh }
	dispmanx.UpdateSubmitSync(dispmanx_update)

	// make surface
	this.surfacep,err = CreateWindowSurface(this.displayp,configs[0],window,nil)
	if err != nil {
		return nil,err
	}

	// connect the context to the surface
	if err := MakeCurrent(this.displayp,this.surfacep,this.surfacep,this.contextp); err != nil {
		return nil,err
	}

	// Success
	return this,nil
}

func (this *EGL) Terminate () error {
	if this.surfacep != EGL_NO_SURFACE {
		// Terminate surface
		if err := DestroySurface(this.displayp,this.surfacep); err != nil {
			return err
		}
	}
	if this.contextp != EGL_NO_CONTEXT {
		// Terminate context
		if err := DestroyContext(this.displayp,this.contextp); err != nil {
			return err
		}
	}
	if this.displayp != EGL_NO_DISPLAY {
		// Terminate display
		if err := Terminate(this.displayp); err != nil {
			return err
		}
	}

	// success
	return nil
}

////////////////////////////////////////////////////////////////////////////////



////////////////////////////////////////////////////////////////////////////////

func (this *EGL) Display () Display {
	return this.displayp
}

func (this *EGL) Context () Context {
	return this.contextp
}

func (this *EGL) Surface () Surface {
	return this.surfacep
}

////////////////////////////////////////////////////////////////////////////////





