/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi
/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2
  #include <EGL/egl.h>
  #include <VG/openvg.h>
*/
import "C"

/*
import (
	"fmt"
	"os"
	"bufio"
)

type VGDraw struct {
	state *EGLState
	window *EGLWindow
}

const (
	VG_CLEAR_COLOR uint16 = 0x1121
)

func (this *EGLState) CreateWindow(api string,frame *Rectangle) (*EGLWindow, error) {
	var err error

	// CREATE WINDOW
	window := new(EGLWindow)

	// CREATE CONTEXT
	window.config, window.context, err = this.createContext(api)
	if err != nil {
		return nil,err
	}

	// CREATE SCREEN ELEMENT
	update, err := this.vc.UpdateBegin()
	if err != nil {
		return nil,err
	}
	source_frame := &Rectangle{}
	source_frame.Set(Point{  0, 0 }, Size{ frame.Size.Width << 16, frame.Size.Height << 16})
	window.element, err = this.vc.AddElement(update, 0, frame, nil, source_frame)
	if err != nil {
		this.destroyContext(window.context)
		return nil,err
	}
	if err := this.vc.UpdateSubmit(update); err != nil {
		this.destroyContext(window.context)
		return nil,err
	}

	// CREATE SURFACE
	nativewindow := &EGLNativeWindow{ window.element.GetHandle(), int(frame.Size.Width), int(frame.Size.Height)}
	window.surface, err = this.createSurface(window.config, nativewindow)
	if err != nil {
		this.destroyContext(window.context)
		return nil,err
	}

	// Attach context to surface
	if err := this.attachContextToSurface(window.context, window.surface); err != nil {
		this.destroyContext(window.context)
		return nil,err
	}

	// Success
	return window,nil
}

func (this *EGLState) CloseWindow(window *EGLWindow) error {
	// Remove surface
	if err := this.destroySurface(window.surface); err != nil {
		return err
	}
	// Remove element
	update, err := this.vc.UpdateBegin()
	if err != nil {
		return err
	}
	if err := this.vc.RemoveElement(update, window.element ); err != nil {
		return err
	}
	if err := this.vc.UpdateSubmit(update); err != nil {
		return err
	}
	if err := this.destroyContext(window.context); err != nil {
		return err
	}

	return nil
}

func (this *EGLState) Flush(window *EGLWindow) error {
	if err := this.swapBuffer(window.surface); err != nil {
		return err
	}
	return nil
}

func (this *EGLWindow) GetFrame() *Rectangle {
	return &Rectangle{ Point{ 0, 0 }, this.element.GetSize() }
}

////////////////////////////////////////////////////////////////////////////////

func (this *EGLState) Draw(window *EGLWindow) (*VGDraw, error) {
	// Create a drawing context
	draw := new(VGDraw)
	draw.state = this
	draw.window = window
	return draw, nil
}

func (this *VGDraw) Flush() error {
	C.vgFlush()
	if err := this.state.Flush(this.window); err != nil {
		return err
	}
	return nil
}

func (this *VGDraw) Clear(frame *Rectangle,color *Color,alpha uint8) error {
	r := C.VGfloat(float32(color.Red) / 255.0)
	g := C.VGfloat(float32(color.Green) / 255.0)
	b := C.VGfloat(float32(color.Blue) / 255.0)
	a := C.VGfloat(float32(alpha) / 255.0)
	clearColor := []C.VGfloat{ r,g,b,a }
	C.vgSetfv(C.VGParamType(VG_CLEAR_COLOR), C.VGint(4), &clearColor[0]);
	C.vgClear(C.VGint(frame.Point.X), C.VGint(frame.Point.Y), C.VGint(frame.Size.Width), C.VGint(frame.Size.Height));
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func (this *EGLState) Do() error {

	// Create a window 1
	window1, err := this.CreateWindow("OpenVG",&Rectangle{ Point{ 100, 100 }, Size{ 200, 200 } })
	if err != nil {
		return err
	}

	// Create a window 2
	window2, err := this.CreateWindow("OpenVG",&Rectangle{ Point{ 400, 250 }, Size{ 100, 100 } })
	if err != nil {
		return err
	}

	// Draw window1
	draw1, err := this.Draw(window1)
	if err != nil {
		return err
	}
	if err := draw1.Clear(window1.GetFrame(),&Color{ 255, 0, 255 },255); err != nil {
		return err
	}
	if err := draw1.Flush(); err != nil {
		return err
	}

	// Draw window2
	draw2, err := this.Draw(window2)
	if err != nil {
		return err
	}
	if err := draw2.Clear(window2.GetFrame(),&Color{ 255, 0, 0 },255); err != nil {
		return err
	}
	if err := draw2.Flush(); err != nil {
		return err
	}

	// Wait...
	fmt.Println("Press any key to exit...")
	reader := bufio.NewReader(os.Stdin)
    reader.ReadString('\n')

	// Close window
	this.CloseWindow(window1)

	return nil
}
*/


