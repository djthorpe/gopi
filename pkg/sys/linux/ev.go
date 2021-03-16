// +build linux

package linux

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO INTERFACE

/*
 #include <linux/input.h>
 static int _EVIOCGNAME(int len)        { return EVIOCGNAME(len); }
 static int _EVIOCGPHYS(int len)        { return EVIOCGPHYS(len); }
 static int _EVIOCGUNIQ(int len)        { return EVIOCGUNIQ(len); }
 static int _EVIOCGPROP(int len)        { return EVIOCGPROP(len); }
 static int _EVIOCGKEY(int len)         { return EVIOCGKEY(len); }
 static int _EVIOCGLED(int len)         { return EVIOCGLED(len); }
 static int _EVIOCGSND(int len)         { return EVIOCGSND(len); }
 static int _EVIOCGSW(int len)          { return EVIOCGSW(len); }
 static int _EVIOCGBIT(int ev, int len) { return EVIOCGBIT(ev, len); }
 static int _EVIOCGABS(int abs)         { return EVIOCGABS(abs); }
 static int _EVIOCSABS(int abs)         { return EVIOCSABS(abs); }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EVType     uint16
	EVKeyCode  uint16
	EVLEDState uint8
)

type EVEvent struct {
	Second      uint32
	Microsecond uint32
	Type        EVType
	Code        EVKeyCode
	Value       uint32
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Internal constants
const (
	MAX_IOCTL_SIZE_BYTES = 256
	EV_DEV               = "/dev/input/event"
	EV_PATH_WILDCARD     = "/sys/class/input/event"
)

var (
	EVIOCGNAME = uintptr(C._EVIOCGNAME(MAX_IOCTL_SIZE_BYTES)) // get device name
	EVIOCGPHYS = uintptr(C._EVIOCGPHYS(MAX_IOCTL_SIZE_BYTES)) // get physical location
	EVIOCGUNIQ = uintptr(C._EVIOCGUNIQ(MAX_IOCTL_SIZE_BYTES)) // get unique identifier
	EVIOCGPROP = uintptr(C._EVIOCGPROP(MAX_IOCTL_SIZE_BYTES)) // get device properties
	EVIOCGID   = uintptr(C.EVIOCGID)                          // get device ID
	EVIOCGLED  = uintptr(C._EVIOCGLED(MAX_IOCTL_SIZE_BYTES))  // get LED states
	EVIOCGKEY  = uintptr(C._EVIOCGLED(MAX_IOCTL_SIZE_BYTES))  // get key states
)

// Event types
// See https://www.kernel.org/doc/Documentation/input/event-codes.txt
const (
	EV_SYN       EVType = 0x0000 // Used as markers to separate events
	EV_KEY       EVType = 0x0001 // Used to describe state changes of keyboards, buttons
	EV_REL       EVType = 0x0002 // Used to describe relative axis value changes
	EV_ABS       EVType = 0x0003 // Used to describe absolute axis value changes
	EV_MSC       EVType = 0x0004 // Miscellaneous uses that didn't fit anywhere else
	EV_SW        EVType = 0x0005 // Used to describe binary state input switches
	EV_LED       EVType = 0x0011 // Used to turn LEDs on devices on and off
	EV_SND       EVType = 0x0012 // Sound output, such as buzzers
	EV_REP       EVType = 0x0014 // Enables autorepeat of keys in the input core
	EV_FF        EVType = 0x0015 // Sends force-feedback effects to a device
	EV_PWR       EVType = 0x0016 // Power management events
	EV_FF_STATUS EVType = 0x0017 // Device reporting of force-feedback effects back to the host
	EV_MAX       EVType = 0x001F
)

const (
	EV_CODE_X        EVKeyCode = 0x0000
	EV_CODE_Y        EVKeyCode = 0x0001
	EV_CODE_SCANCODE EVKeyCode = 0x0004 // Keyboard scan code
	EV_CODE_SLOT     EVKeyCode = 0x002F // Slot for multi touch positon
	EV_CODE_SLOT_X   EVKeyCode = 0x0035 // X for multi touch position
	EV_CODE_SLOT_Y   EVKeyCode = 0x0036 // Y for multi touch position
	EV_CODE_SLOT_ID  EVKeyCode = 0x0039 // Unique ID for multi touch position
)

// LED Constants
const (
	EV_LED_NUML     EVLEDState = 0x00
	EV_LED_CAPSL    EVLEDState = 0x01
	EV_LED_SCROLLL  EVLEDState = 0x02
	EV_LED_COMPOSE  EVLEDState = 0x03
	EV_LED_KANA     EVLEDState = 0x04
	EV_LED_SLEEP    EVLEDState = 0x05
	EV_LED_SUSPEND  EVLEDState = 0x06
	EV_LED_MUTE     EVLEDState = 0x07
	EV_LED_MISC     EVLEDState = 0x08
	EV_LED_MAIL     EVLEDState = 0x09
	EV_LED_CHARGING EVLEDState = 0x0A
	EV_LED_MAX      EVLEDState = 0x0F
)

////////////////////////////////////////////////////////////////////////////////
// OPEN

func EVDevice(bus uint) string {
	return fmt.Sprintf("%v%v", EV_DEV, bus)
}

func EVDevices() ([]uint, error) {
	if files, err := filepath.Glob(EV_PATH_WILDCARD + "*"); err != nil {
		return nil, err
	} else {
		devices := make([]uint, 0, len(files))
		for _, file := range files {
			if bus, err := strconv.ParseUint(strings.TrimPrefix(file, EV_PATH_WILDCARD), 10, 32); err == nil {
				devices = append(devices, uint(bus))
			}
		}
		return devices, nil
	}
}

func EVOpenDevice(bus uint) (*os.File, error) {
	if file, err := os.OpenFile(EVDevice(bus), os.O_SYNC|os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// IOCTL FUNCTIONS

// Get name of the device
func EVGetName(fd uintptr) (string, error) {
	str := make([]C.char, MAX_IOCTL_SIZE_BYTES)
	if err := ev_ioctl(fd, uintptr(EVIOCGNAME), unsafe.Pointer(&str[0])); err != nil {
		return "", err
	} else {
		return C.GoString(&str[0]), nil
	}
}

// Get physical connection string
func EVGetPhys(fd uintptr) (string, error) {
	name := new([MAX_IOCTL_SIZE_BYTES]C.char)
	if err := ev_ioctl(fd, uintptr(EVIOCGPHYS), unsafe.Pointer(name)); err != nil {
		return "", err
	} else {
		return C.GoString(&name[0]), nil
	}
}

// Get unique identifier
func EVGetUniq(fd uintptr) (string, error) {
	name := new([MAX_IOCTL_SIZE_BYTES]C.char)
	if err := ev_ioctl(fd, uintptr(EVIOCGUNIQ), unsafe.Pointer(name)); err != nil {
		return "", err
	} else {
		return C.GoString(&name[0]), nil
	}
}

// Get device information (bus, vendor, product, version)
func EVGetInfo(fd uintptr) ([]uint16, error) {
	info := make([]uint16, 4)
	if err := ev_ioctl(fd, uintptr(EVIOCGID), unsafe.Pointer(&info[0])); err != nil {
		return nil, err
	} else {
		return info, nil
	}
}

// Get device capabilities
func EVGetSupportedEventTypes(fd uintptr) ([]EVType, error) {
	evbits := new([EV_MAX >> 3]byte)
	if err := ev_ioctl(fd, uintptr(C._EVIOCGBIT(C.int(0), C.int(EV_MAX))), unsafe.Pointer(evbits)); err != nil {
		return nil, err
	}
	capabilities := make([]EVType, 0, EV_MAX)
	evtype := EVType(0)
	for i := 0; i < len(evbits); i++ {
		evbyte := evbits[i]
		for j := 0; j < 8; j++ {
			if evbyte&0x01 != 0x00 {
				capabilities = append(capabilities, evtype)
			}
			evbyte = evbyte >> 1
			evtype++
		}
	}
	return capabilities, nil
}

// Obtain and release exclusive device usage ("grab")
func EVSetGrabState(fd uintptr, state bool) error {
	if state {
		if err := ev_ioctl(fd, C.EVIOCGRAB, unsafe.Pointer(uintptr(1))); err != nil {
			return err
		}
	} else {
		if err := ev_ioctl(fd, C.EVIOCGRAB, unsafe.Pointer(uintptr(0))); err != nil {
			return err
		}
	}
	return nil
}

// EVGetLEDState gets LED states as an array of LED's which are on
func EVGetLEDState(fd uintptr) ([]EVLEDState, error) {
	evbits := new([MAX_IOCTL_SIZE_BYTES]byte)
	if err := ev_ioctl(fd, EVIOCGLED, unsafe.Pointer(evbits)); err != nil {
		return nil, err
	}
	// Shift bits to get the state of each LED value
	states := make([]EVLEDState, 0, EV_LED_MAX)
FOR_LOOP:
	for i := 0; i < len(evbits); i++ {
		evbyte := evbits[i]
		for j := 0; j < 8; j++ {
			state := EVLEDState(i<<3 + j)
			switch {
			case state >= EV_LED_MAX:
				break FOR_LOOP
			case evbyte&0x01 != 0x00:
				states = append(states, state)
			}
			evbyte >>= 1
		}
	}
	return states, nil
}

// EVSetLEDState sets a single LED state
func EVSetLEDState(fd uintptr, led EVLEDState, state bool) error {
	file := os.NewFile(fd, "EVSetLEDState")
	if state {
		return binary.Write(file, binary.LittleEndian, EVEvent{
			Type:  EV_LED,
			Code:  EVKeyCode(led),
			Value: 1,
		})
	} else {
		return binary.Write(file, binary.LittleEndian, EVEvent{
			Type:  EV_LED,
			Code:  EVKeyCode(led),
			Value: 0,
		})
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e EVEvent) String() string {
	return "<ev_event" +
		" type=" + fmt.Sprint(e.Type) +
		" code=" + fmt.Sprint(e.Code) +
		" value=" + fmt.Sprintf("0x%08X", e.Value) +
		">"
}

func (v EVType) String() string {
	switch v {
	case EV_SYN:
		return "EV_SYN"
	case EV_KEY:
		return "EV_KEY"
	case EV_REL:
		return "EV_REL"
	case EV_ABS:
		return "EV_ABS"
	case EV_MSC:
		return "EV_MSC"
	case EV_SW:
		return "EV_SW"
	case EV_LED:
		return "EV_LED"
	case EV_SND:
		return "EV_SND"
	case EV_REP:
		return "EV_REP"
	case EV_FF:
		return "EV_FF"
	case EV_PWR:
		return "EV_PWR"
	case EV_FF_STATUS:
		return "EV_FF_STATUS"
	default:
		return "[?? Invalid EVType value]"
	}
}

func (v EVKeyCode) String() string {
	switch v {
	case EV_CODE_X:
		return "EV_CODE_X"
	case EV_CODE_Y:
		return "EV_CODE_Y"
	case EV_CODE_SCANCODE:
		return "EV_CODE_SCANCODE"
	case EV_CODE_SLOT:
		return "EV_CODE_SLOT"
	case EV_CODE_SLOT_X:
		return "EV_CODE_SLOT_X"
	case EV_CODE_SLOT_Y:
		return "EV_CODE_SLOT_Y"
	case EV_CODE_SLOT_ID:
		return "EV_CODE_SLOT_ID"
	default:
		return "[?? Invalid EVKeyCode value]"
	}
}

func (s EVLEDState) String() string {
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
		return "[?? Invalid EVLEDState value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Call ioctl
func ev_ioctl(fd uintptr, name uintptr, data unsafe.Pointer) error {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	if err != 0 {
		return os.NewSyscallError("ev_ioctl", err)
	} else {
		return nil
	}
}
