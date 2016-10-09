/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

import (
	rpi "../device/rpi"
)

var (
	flagDisplay = flag.String("display", "lcd", "Display")
)

func CreateBackground(vc *rpi.VideoCore,color rpi.Color,alpha uint8) (*rpi.Resource,error) {
	// create a 1x1 resource
	resource, err := vc.CreateResource(rpi.VC_IMAGE_RGBA32,rpi.Size{ 1, 1 })
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

func CreateGradient(vc *rpi.VideoCore, nw,ne,se,sw rpi.Color,alpha uint8) (*rpi.Resource,error) {
	// create a 2x2 resource
	resource, err := vc.CreateResource(rpi.VC_IMAGE_RGBA32,rpi.Size{ 2, 2 })
	if err != nil {
		return nil, err
	}
	// write data into the resource
	pixel := make([]byte,128) // two lines of 16 x uint32

	// nw pixel
	pixel[0] = nw.Red
	pixel[1] = nw.Green
	pixel[2] = nw.Blue
	pixel[3] = alpha

	// ne pixel
	pixel[4] = ne.Red
	pixel[5] = ne.Green
	pixel[6] = ne.Blue
	pixel[7] = alpha

	// sw pixel - start at 16*4 = 64
	pixel[64] = sw.Red
	pixel[65] = sw.Green
	pixel[66] = sw.Blue
	pixel[67] = alpha

	// se pixel
	pixel[68] = se.Red
	pixel[69] = se.Green
	pixel[70] = se.Blue
	pixel[71] = alpha

	// write the data
	// uint32 is 4 bytes and one row is 16 uint32 - so stride is 16 * 4
	err = resource.WriteData(rpi.VC_IMAGE_RGBA32,16 * 4,pixel,resource.GetFrame())
	if err != nil {
		return nil, err
	}
	// success
	return resource, nil
}


func main() {
	// Flags
	flag.Parse()

	// Retrieve display
	display, ok := rpi.Displays[*flagDisplay]
	if !ok {
		fmt.Fprintln(os.Stderr, "Error: Invalid display name")
		os.Exit(-1)
	}

	// Open up the RaspberryPi interface
	pi, err := rpi.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	defer pi.Close()

	// VideoCore
	vc, err := pi.NewVideoCore(display)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	defer vc.Close()

	// Create a background & two gradients
	bg, _ := CreateBackground(vc,rpi.Color{ 70, 70, 70 },255)
	fg1, _ := CreateGradient(vc,rpi.Color{ 255, 0, 0 },rpi.Color{ 255, 255, 0 },rpi.Color{ 0, 255, 255 },rpi.Color{ 255, 0, 255 },255)
	fg2, _ := CreateGradient(vc,rpi.Color{ 255, 255, 0 },rpi.Color{ 0, 255, 255 },rpi.Color{ 255, 0, 255 },rpi.Color{ 255, 0, 0 },255)

	// Start an update
	update, err := vc.UpdateBegin()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}

	// Add background element stretch to fill screen
	vc.AddElement(update,1,&rpi.Rectangle{},bg,bg.GetFrame())

	// Add foreground element to overlay
	size := vc.GetSize()
	position1 := &rpi.Rectangle{ rpi.Point{ 100,100 }, rpi.Size{ size.Width - 200, size.Height - 200 } }
	position2 := &rpi.Rectangle{ rpi.Point{ 150,150 }, rpi.Size{ size.Width - 200, size.Height - 200 } }
	element1, _ := vc.AddElement(update,2,position1,fg1,&rpi.Rectangle{ rpi.Point{}, rpi.Size{ 2 << 16, 2 << 16 }});
	element2, _ := vc.AddElement(update,3,position2,fg2,&rpi.Rectangle{ rpi.Point{}, rpi.Size{ 2 << 16, 2 << 16 }});

	// Place elements on the screen
	vc.UpdateSubmit(update)

	// now switch between the two sources
	for {
		if position1.Size.Width <= 0 || position1.Size.Height <= 0 {
			break
		}
		if position1.Point.X <= 0 || position1.Point.Y <= 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)

		update, _ := vc.UpdateBegin()

		position1.Point.X -= 1
		position1.Point.Y -= 1
		vc.ChangeElementFrame(update,element1,position1)

		position2.Point.X += 3
		position2.Point.Y += 3
		vc.ChangeElementFrame(update,element2,position2)

		vc.UpdateSubmit(update)
	}

	// Remove the resources
	vc.DeleteResource(fg1)
	vc.DeleteResource(fg2)
	vc.DeleteResource(bg)
}
