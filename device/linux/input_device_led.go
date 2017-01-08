/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"encoding/binary"
	"os"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type evLEDState uint8

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// LED Constants
const (
	EV_LED_NUML     evLEDState = 0x00
	EV_LED_CAPSL    evLEDState = 0x01
	EV_LED_SCROLLL  evLEDState = 0x02
	EV_LED_COMPOSE  evLEDState = 0x03
	EV_LED_KANA     evLEDState = 0x04
	EV_LED_SLEEP    evLEDState = 0x05
	EV_LED_SUSPEND  evLEDState = 0x06
	EV_LED_MUTE     evLEDState = 0x07
	EV_LED_MISC     evLEDState = 0x08
	EV_LED_MAIL     evLEDState = 0x09
	EV_LED_CHARGING evLEDState = 0x0A
	EV_LED_MAX      evLEDState = 0x0F
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s evLEDState) String() string {
	switch s {
	case EV_LED_NUML:
		return "EV_LED_NUML"
	case EV_LED_CAPSL:
		return "EV_LED_CAPSL"
	case EV_LED_SCROLLL:
		return "EV_LED_SCROLLL"
	case EV_LED_COMPOSE:
		return "EV_LED_COMPOSE"
	case EV_LED_KANA:
		return "EV_LED_KANA"
	case EV_LED_SLEEP:
		return "EV_LED_SLEEP"
	case EV_LED_SUSPEND:
		return "EV_LED_SUSPEND"
	case EV_LED_MUTE:
		return "EV_LED_MUTE"
	case EV_LED_MISC:
		return "EV_LED_MISC"
	case EV_LED_MAIL:
		return "EV_LED_MAIL"
	case EV_LED_CHARGING:
		return "EV_LED_CHARGING"
	default:
		return "[?? Invalid evLEDState value]"
	}
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Get LED states as an array of LED's which are on
func evGetLEDState(handle *os.File) ([]evLEDState, error) {
	evbits := new([MAX_IOCTL_SIZE_BYTES]byte)
	err := evIoctl(handle.Fd(), uintptr(EVIOCGLED), unsafe.Pointer(evbits))
	if err != 0 {
		return nil, err
	}
	states := make([]evLEDState, 0, EV_LED_MAX)
	// Shift bits to get the state of each LED value
OuterLoop:
	for i := 0; i < len(evbits); i++ {
		evbyte := evbits[i]
		for j := 0; j < 8; j++ {
			state := evLEDState(i<<3 + j)
			switch {
			case state >= EV_LED_MAX:
				break OuterLoop
			case evbyte&0x01 != 0x00:
				states = append(states, state)
			}
			evbyte >>= 1
		}
	}
	return states, nil
}

// Set a single LED state
func evSetLEDState(handle *os.File, led evLEDState, state bool) error {
	var event evEvent

	event.Type = EV_LED
	event.Code = evKeyCode(led)

	if state {
		event.Value = 1
	}
	if err := binary.Write(handle, binary.LittleEndian, &event); err != nil {
		return err
	}
	return nil
}
