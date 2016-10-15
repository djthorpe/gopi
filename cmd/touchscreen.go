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
)

import (
	"../input"
	"../device/touchscreen/ft5406"
)

func main() {
	touchscreen, err := input.Open(ft5406.Config{})
	if err != nil {
		fmt.Println("Error: ",err)
		os.Exit(-1)
	}
	defer touchscreen.Close()

	fmt.Println("Device:",touchscreen.GetName())

	err = touchscreen.ProcessTouchEvents(func(dev *input.Device, evt *input.TouchEvent) {
		switch {
		case evt.Type==input.EVENT_BTN_PRESS:
			//fmt.Println("PRESS:")
			break
		case evt.Type==input.EVENT_BTN_RELEASE:
			//fmt.Println("RELEASE:")
			break
		case evt.Type==input.EVENT_MOVE:
			//fmt.Println("MOVE:",evt.LastPoint,"->",evt.Point)
			break
		case evt.Type==input.EVENT_SLOT_MOVE:
			fmt.Println("SLOT MOVE:",evt.Slot,evt.LastPoint,"->",evt.Point)
			break
		case evt.Type==input.EVENT_SLOT_PRESS:
			fmt.Println("SLOT PRESS:",evt.Slot,evt.Point)
			break
		case evt.Type==input.EVENT_SLOT_RELEASE:
			fmt.Println("SLOT RELEASE:",evt.Slot,evt.Point)
			break
		}
	})
	if err != nil {
		fmt.Println("Error: ",err)
		os.Exit(-1)
	}
}
