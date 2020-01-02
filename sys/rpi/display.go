// +build rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	DXDisplayId uint16
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// DX_DisplayId values
	DX_DISPLAYID_MAIN_LCD DXDisplayId = iota
	DX_DISPLAYID_AUX_LCD
	DX_DISPLAYID_HDMI
	DX_DISPLAYID_SDTV
	DX_DISPLAYID_FORCE_LCD
	DX_DISPLAYID_FORCE_TV
	DX_DISPLAYID_FORCE_OTHER
	DX_DISPLAYID_MAX = DX_DISPLAYID_FORCE_OTHER
	DX_DISPLAYID_MIN = DX_DISPLAYID_MAIN_LCD
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func DXNumberOfDisplays() uint16 {
	return uint16(DX_DISPLAYID_MAX-DX_DISPLAYID_MIN) + 1
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (d DXDisplayId) String() string {
	switch d {
	case DX_DISPLAYID_MAIN_LCD:
		return "DX_DISPLAYID_MAIN_LCD"
	case DX_DISPLAYID_AUX_LCD:
		return "DX_DISPLAYID_AUX_LCD"
	case DX_DISPLAYID_HDMI:
		return "DX_DISPLAYID_HDMI"
	case DX_DISPLAYID_SDTV:
		return "DX_DISPLAYID_SDTV"
	case DX_DISPLAYID_FORCE_LCD:
		return "DX_DISPLAYID_FORCE_LCD"
	case DX_DISPLAYID_FORCE_TV:
		return "DX_DISPLAYID_FORCE_TV"
	case DX_DISPLAYID_FORCE_OTHER:
		return "DX_DISPLAYID_FORCE_OTHER"
	default:
		return "[?? Invalid DXDisplayId value]"
	}
}
