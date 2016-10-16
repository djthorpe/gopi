/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"os"
	"fmt"
	"flag"
	"image"
)

import (
	input "../input"
	ft5406 "../device/touchscreen/ft5406"
	rpi "../device/rpi"
)

////////////////////////////////////////////////////////////////////////////////

type Application struct {
	rpi *rpi.RaspberryPi
	touchscreen *input.Device
	vc *rpi.VideoCore
}

////////////////////////////////////////////////////////////////////////////////

var (
	flagDisplay = flag.String("display", "lcd", "Display")
)

////////////////////////////////////////////////////////////////////////////////

func NewApplication(display uint16) (*Application, error) {
	var err error

	app := new(Application)

	// Raspberry Pi
	app.rpi,err = rpi.New()
	if err != nil {
		return nil, err
	}

	// Touchscreen
	app.touchscreen, err = input.Open(ft5406.Config{})
	if err != nil {
		app.rpi.Close()
		return nil,err
	}

	// VideoCore
	app.vc, err = app.rpi.NewVideoCore(display)
	if err != nil {
		app.rpi.Close()
		app.touchscreen.Close()
		return nil, err
	}

	return app, nil
}

func (app *Application) Close() error {
	app.vc.Close()
	app.touchscreen.Close()
	app.rpi.Close()
	return nil
}

func (app *Application) CreatePixel(color rpi.Color,alpha uint8) (*rpi.Resource,error) {
	// create a 1x1 resource
	resource, err := app.vc.CreateResource(rpi.VC_IMAGE_RGBA32,rpi.Size{ 1, 1 })
	if err != nil {
		return nil, err
	}
	// write data into the resource
	var pixel []byte
	pixel = make([]byte,4)
	pixel[0] = color.Red
	pixel[1] = color.Green
	pixel[2] = color.Blue
	pixel[3] = alpha
	err = resource.WriteData(rpi.VC_IMAGE_RGBA32,len(pixel),pixel,resource.GetFrame())
	if err != nil {
		return nil, err
	}
	// success
	return resource, nil
}

func (app *Application) AddElement(slot uint32,point image.Point) (*rpi.Element, error) {
	width := uint32(50)
	height := uint32(50)

	resource, err := app.CreatePixel(rpi.Color{ 255, 0, 0 },255)
	if err != nil {
		return nil, err
	}
	// Start an update
	update, err := app.vc.UpdateBegin()
	if err != nil {
		return nil, err
	}

	// Add background element stretch to fill screen
	frame := &rpi.Rectangle{ rpi.Point{ int32(point.X) - int32(width/2), int32(point.Y) - int32(height/2) }, rpi.Size{ width, height } }
	element, err := app.vc.AddElement(update,2,frame,resource,&rpi.Rectangle{ rpi.Point{}, rpi.Size{ 1 << 16, 1 << 16 }});
	if err != nil {
		app.vc.DeleteResource(resource)
		return nil, err
	}
	app.vc.UpdateSubmit(update)
	return element, nil
}

func (app *Application) DeleteElement(element *rpi.Element) error {
	// TODO
	return nil
}

func (app *Application) MoveElement(element *rpi.Element,point image.Point) error {
	// Start an update
	update, err := app.vc.UpdateBegin()
	if err != nil {
		return err
	}

	rect := element.GetFrame()
	rect.X = int32(point.X) - int32(rect.Width/2)
	rect.Y = int32(point.Y) - int32(rect.Height/2)
	app.vc.ChangeElementFrame(update,element,rect)
	app.vc.UpdateSubmit(update)
	return nil
}

func (app *Application) Run() error {

	// Create a background window
	bg, err := app.CreatePixel(rpi.Color{ 70, 70, 190 },255)
	if err != nil {
		return err
	}
	defer app.vc.DeleteResource(bg)

	// Start an update
	update, err := app.vc.UpdateBegin()
	if err != nil {
		return err
	}

	// Add background element stretch to fill screen
	app.vc.AddElement(update,1,&rpi.Rectangle{},bg,bg.GetFrame())
	app.vc.UpdateSubmit(update)

	// Set up slots
	slots := make(map[uint32]*rpi.Element)

	err = app.touchscreen.ProcessTouchEvents(func(dev *input.Device, evt *input.TouchEvent) {
		switch {
		case evt.Type==input.EVENT_SLOT_PRESS:
			if slots[evt.Slot] != nil {
				app.DeleteElement(slots[evt.Slot])
			}
			slots[evt.Slot], _ = app.AddElement(evt.Slot,evt.Point)
			break
		case evt.Type==input.EVENT_SLOT_MOVE:
			if slots[evt.Slot] != nil {
				app.MoveElement(slots[evt.Slot],evt.Point)
			}
			break
		case evt.Type==input.EVENT_SLOT_RELEASE:
			if slots[evt.Slot] != nil {
				app.DeleteElement(slots[evt.Slot])
			}
			break
		}
	})

	return err
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Flags
	flag.Parse()

	// Retrieve display
	display, ok := rpi.Displays[*flagDisplay]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: Invalid display name")
		os.Exit(-1)
	}

	// Create an application object
	app, err := NewApplication(display)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(-1)
	}
	defer app.Close()

	err = app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

