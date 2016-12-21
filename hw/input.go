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

	// Get the bus interface
	GetBus() InputDeviceBus

	// Close the device
	Close() error
}

// Device type (keyboard, mouse, touchscreen, etc)
type InputDeviceType uint8

// Bus type (USB, Bluetooth, etc)
type InputDeviceBus uint16

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

// Types of input connection
const (
	INPUT_BUS_NONE      InputDeviceBus = 0x0000
	INPUT_BUS_PCI       InputDeviceBus = 0x0001
	INPUT_BUS_ISAPNP    InputDeviceBus = 0x0002
	INPUT_BUS_USB       InputDeviceBus = 0x0003
	INPUT_BUS_HIL       InputDeviceBus = 0x0004
	INPUT_BUS_BLUETOOTH InputDeviceBus = 0x0005
	INPUT_BUS_VIRTUAL   InputDeviceBus = 0x0006
	INPUT_BUS_ISA       InputDeviceBus = 0x0010
	INPUT_BUS_I8042     InputDeviceBus = 0x0011
	INPUT_BUS_XTKBD     InputDeviceBus = 0x0012
	INPUT_BUS_RS232     InputDeviceBus = 0x0013
	INPUT_BUS_GAMEPORT  InputDeviceBus = 0x0014
	INPUT_BUS_PARPORT   InputDeviceBus = 0x0015
	INPUT_BUS_AMIGA     InputDeviceBus = 0x0016
	INPUT_BUS_ADB       InputDeviceBus = 0x0017
	INPUT_BUS_I2C       InputDeviceBus = 0x0018
	INPUT_BUS_HOST      InputDeviceBus = 0x0019
	INPUT_BUS_GSC       InputDeviceBus = 0x001A
	INPUT_BUS_ATARI     InputDeviceBus = 0x001B
	INPUT_BUS_SPI       InputDeviceBus = 0x001C
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY FUNCTIONS

func (t InputDeviceType) String() string {
	switch t {
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

func (b InputDeviceBus) String() string {
	switch b {
	case INPUT_BUS_NONE:
		return "INPUT_BUS_NONE"
	case INPUT_BUS_PCI:
		return "INPUT_BUS_PCI"
	case INPUT_BUS_ISAPNP:
		return "INPUT_BUS_ISAPNP"
	case INPUT_BUS_USB:
		return "INPUT_BUS_USB"
	case INPUT_BUS_HIL:
		return "INPUT_BUS_HIL"
	case INPUT_BUS_BLUETOOTH:
		return "INPUT_BUS_BLUETOOTH"
	case INPUT_BUS_VIRTUAL:
		return "INPUT_BUS_VIRTUAL"
	case INPUT_BUS_ISA:
		return "INPUT_BUS_ISA"
	case INPUT_BUS_I8042:
		return "INPUT_BUS_I8042"
	case INPUT_BUS_XTKBD:
		return "INPUT_BUS_XTKBD"
	case INPUT_BUS_RS232:
		return "INPUT_BUS_RS232"
	case INPUT_BUS_GAMEPORT:
		return "INPUT_BUS_GAMEPORT"
	case INPUT_BUS_PARPORT:
		return "INPUT_BUS_PARPORT"
	case INPUT_BUS_AMIGA:
		return "INPUT_BUS_AMIGA"
	case INPUT_BUS_ADB:
		return "INPUT_BUS_ADB"
	case INPUT_BUS_I2C:
		return "INPUT_BUS_I2C"
	case INPUT_BUS_HOST:
		return "INPUT_BUS_HOST"
	case INPUT_BUS_GSC:
		return "INPUT_BUS_GSC"
	case INPUT_BUS_ATARI:
		return "INPUT_BUS_ATARI"
	case INPUT_BUS_SPI:
		return "INPUT_BUS_SPI"
	default:
		return "[?? Invalid InputDeviceBus value]"
	}
}
