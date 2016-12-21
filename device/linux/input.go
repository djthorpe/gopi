/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
	util "github.com/djthorpe/gopi/util"
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

// Empty input configuration
type Input struct{}

// Driver of multiple input devices
type InputDriver struct {
	log     *util.LoggerDevice // logger
	devices []hw.InputDevice   // input devices
}

// A single input device
type InputDevice struct {
	// The name of the input device
	Name string

	// The device path to the input device
	Path string

	// The Id of the input device
	Id string

	// The type of device, or NONE
	Type hw.InputDeviceType

	// The bus which the device is attached to, or NONE
	Bus hw.InputDeviceBus

	// Product and version
	Vendor uint16
	Product uint16
	Version uint16

	// Capabilities
	Events []evType

	// Handle to the device
	handle *os.File
}

type evType uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Internal constants
const (
	PATH_INPUT_DEVICES   = "/sys/class/input/event*"
	MAX_POLL_EVENTS      = 32
	MAX_EVENT_SIZE_BYTES = 1024
	MAX_IOCTL_SIZE_BYTES = 256
)

// Event types
// See https://www.kernel.org/doc/Documentation/input/event-codes.txt
const (
	EV_SYN       evType = 0x0000 // Used as markers to separate events
	EV_KEY       evType = 0x0001 // Used to describe state changes of keyboards, buttons
	EV_REL       evType = 0x0002 // Used to describe relative axis value changes
	EV_ABS       evType = 0x0003 // Used to describe absolute axis value changes
	EV_MSC       evType = 0x0004 // Miscellaneous uses that didn't fit anywhere else
	EV_SW        evType = 0x0005 // Used to describe binary state input switches
	EV_LED       evType = 0x0011 // Used to turn LEDs on devices on and off
	EV_SND       evType = 0x0012 // Sound output, such as buzzers
	EV_REP       evType = 0x0014 // Enables autorepeat of keys in the input core
	EV_FF        evType = 0x0015 // Sends force-feedback effects to a device
	EV_PWR       evType = 0x0016 // Power management events
	EV_FF_STATUS evType = 0x0017 // Device reporting of force-feedback effects back to the host
	EV_MAX       evType = 0x001F
)

var (
	EVIOCGNAME = uintptr(C._EVIOCGNAME(MAX_IOCTL_SIZE_BYTES)) // get device name
	EVIOCGPHYS = uintptr(C._EVIOCGPHYS(MAX_IOCTL_SIZE_BYTES)) // get physical location
	EVIOCGUNIQ = uintptr(C._EVIOCGUNIQ(MAX_IOCTL_SIZE_BYTES)) // get unique identifier
	EVIOCGPROP = uintptr(C._EVIOCGPROP(MAX_IOCTL_SIZE_BYTES)) // get device properties
	EVIOCGID   = uintptr(C.EVIOCGID)                          // get device ID
)

////////////////////////////////////////////////////////////////////////////////
// InputDriver OPEN AND CLOSE

// Create new Input object, returns error if not possible
func (config Input) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<linux.Input>Open")

	// create new GPIO driver
	this := new(InputDriver)

	// Set logging & device
	this.log = log

	// Find devices
	this.devices = make([]hw.InputDevice, 0)
	if err := evFind(func(device *InputDevice) {
		this.devices = append(this.devices, device)
	}); err != nil {
		return nil, err
	}

	// Get capabilities for devices
	for _, device := range this.devices {
		err := device.(*InputDevice).Open()
		defer device.(*InputDevice).Close()
		if err == nil {
			err = device.(*InputDevice).evSetCapabilities()
		}
		if err != nil {
			log.Warn("Device %v: %v",device.GetName(),err)
		}
	}

	// success
	return this, nil
}

// Close Input driver
func (this *InputDriver) Close() error {
	this.log.Debug("<linux.Input>Close")

	for _, device := range this.devices {
		if err := device.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (t evType) String() string {
	switch(t) {
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
		return "[?? Unknown evType value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// InputDevice OPEN AND CLOSE

// Open driver
func (this *InputDevice) Open() error {
	if this.handle != nil {
		if err := this.Close(); err != nil {
			return err
		}
	}
	var err error
	if this.handle, err = os.OpenFile(this.Path, os.O_RDWR, 0); err != nil {
		this.handle = nil
		return err
	}
	// Success
	return nil
}

// Close driver
func (this *InputDevice) Close() error {
	var err error
	if this.handle != nil {
		err = this.handle.Close()
	}
	this.handle = nil
	return err
}

////////////////////////////////////////////////////////////////////////////////
// InputDevice implementation

func (this *InputDevice) GetName() string {
	return this.Name
}

func (this *InputDevice) GetType() hw.InputDeviceType {
	return this.Type
}

func (this *InputDevice) GetBus() hw.InputDeviceBus {
	return this.Bus
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

// Strinfigy InputDriver object
func (this *InputDriver) String() string {
	return fmt.Sprintf("<linux.Input>{ devices=%v }", this.devices)
}

// Strinfigy InputDevice object
func (this *InputDevice) String() string {
	return fmt.Sprintf("<linux.InputDevice>{ name=\"%s\" path=%s id=%v type=%v bus=%v product=0x%04X vendor=0x%04X version=0x%04X events=%v fd=%v }", this.Name, this.Path, this.Id, this.Type, this.Bus, this.Product, this.Vendor, this.Version, this.Events, this.handle)
}



////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *InputDevice) evSetCapabilities() error {
	name, err := evGetName(this.handle)
	if err != nil {
		return err
	}
	this.Name = name

	id, err := evGetPhys(this.handle)
	// Error is ignored
	if err == nil {
		this.Id = id
	}

	bus, vendor, product, version, err := evGetInfo(this.handle)
	if err == nil {
		// Error is ignored
		this.Bus = hw.InputDeviceBus(bus)
		this.Vendor = vendor
		this.Product = product
		this.Version = version
	}

	events, err := evGetEvents(this.handle)
	if err != nil {
		return err
	}
	this.Events = events

	return nil
}

// Find all input devices
func evFind(callback func(driver *InputDevice)) error {
	files, err := filepath.Glob(PATH_INPUT_DEVICES)
	if err != nil {
		return err
	}
	for _, file := range files {
		buf, err := ioutil.ReadFile(path.Join(file, "device", "name"))
		if err != nil {
			continue
		}
		device := &InputDevice{Name: strings.TrimSpace(string(buf)), Path: path.Join("/", "dev", "input", path.Base(file))}
		callback(device)
	}
	return nil
}

// Get name
func evGetName(handle *os.File) (string, error) {
	name := new([MAX_IOCTL_SIZE_BYTES]C.char)
	err := evIoctl(handle.Fd(), uintptr(EVIOCGNAME), unsafe.Pointer(name))
	if err != 0 {
		return "", err
	}
	return C.GoString(&name[0]), nil
}

// Get physical connection string
func evGetPhys(handle *os.File) (string, error) {
	name := new([MAX_IOCTL_SIZE_BYTES]C.char)
	err := evIoctl(handle.Fd(), uintptr(EVIOCGPHYS), unsafe.Pointer(name))
	if err != 0 {
		return "", err
	}
	return C.GoString(&name[0]), nil
}

// Get device information (bus, vendor, product, version)
func evGetInfo(handle *os.File) (uint16,uint16,uint16,uint16,error) {
	info := [4]uint16{ }
	err := evIoctl(handle.Fd(), uintptr(EVIOCGID), unsafe.Pointer(&info))
	if err != 0 {
		return uint16(0),uint16(0),uint16(0),uint16(0),err
	}
	return info[0],info[1],info[2],info[3],nil
}

// Get supported events
func evGetEvents(handle *os.File) ([]evType,error) {
	evbits := new([EV_MAX >> 3]byte)
	err := evIoctl(handle.Fd(),uintptr(C._EVIOCGBIT(C.int(0), C.int(EV_MAX))), unsafe.Pointer(evbits))
	if err != 0 {
		return nil,err
	}
	capabilities := make([]evType,0)
	evtype := evType(0)
	for i := 0; i < len(evbits); i++ {
		evbyte := evbits[i]
		for j := 0; j < 8; j++ {
			if evbyte & 0x01 != 0x00 {
				capabilities = append(capabilities,evtype)
			}
			evbyte = evbyte >> 1
			evtype++
		}
	}
	return capabilities,nil
}

// Call ioctl
func evIoctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}

