/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package hw // import "github.com/djthorpe/gopi/hw"

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type InputDriver interface {
	// Enforces general driver
	gopi.Driver

	// Return devices
}

type InputDevice interface {
	// Get the name of the input device
	GetName() string

	// Get the type of device
	GetType() InputDeviceType

	// Close the device
	Close() error
}

type InputDeviceType uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Types of input device
const (
	INPUT_TYPE_NONE        InputDeviceType = 0x00
	INPUT_TYPE_KEYBOARD    InputDeviceType = 0x01
	INPUT_TYPE_MOUSE       InputDeviceType = 0x02
	INPUT_TYPE_TOUCHSCREEN InputDeviceType = 0x04
	INPUT_TYPE_JOYSTICK    InputDeviceType = 0x08
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY FUNCTIONS

func (t InputDeviceType) String() string {
	switch(t) {
	case INPUT_TYPE_NONE:
		return "INPUT_TYPE_NONE"
	case INPUT_TYPE_KEYBOARD:
		return "INPUT_TYPE_KEYBOARD"
	case INPUT_TYPE_MOUSE:
		return "INPUT_TYPE_MOUSE"
	case INPUT_TYPE_TOUCHSCREEN:
		return "INPUT_TYPE_TOUCHSCREEN"
	case INPUT_TYPE_JOYSTICK:
		return "INPUT_TYPE_JOYSTICK"
	default:
		return "[?? Invalid InputDeviceType value]"
	}
}

