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
	"path"
	"path/filepath"
	"strings"
	"os"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Empty input configuration
type Input struct{}

type InputDriver struct {
	log *util.LoggerDevice // logger
	devices []hw.InputDevice // input devices
}

type InputDevice struct {
	// The name of the input device
	Name string

	// The path to the input device
	Path string

	// The type of device, or NONE
	Type hw.InputDeviceType

	// Handle to the device
	handle *os.File
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

// Internal constants
const (
	PATH_INPUT_DEVICES = "/sys/class/input/event*"
	MAX_POLL_EVENTS      = 32
	MAX_EVENT_SIZE_BYTES = 1024
)

// Event types
// See https://www.kernel.org/doc/Documentation/input/event-codes.txt
const (
	EV_SYN       uint16 = 0x0000 // Used as markers to separate events
	EV_KEY       uint16 = 0x0001 // Used to describe state changes of keyboards, buttons
	EV_REL       uint16 = 0x0002 // Used to describe relative axis value changes
	EV_ABS       uint16 = 0x0003 // Used to describe absolute axis value changes
	EV_MSC       uint16 = 0x0004 // Miscellaneous uses that didn't fit anywhere else
	EV_SW        uint16 = 0x0005 // Used to describe binary state input switches
	EV_LED       uint16 = 0x0011 // Used to turn LEDs on devices on and off
	EV_SND       uint16 = 0x0012 // Sound output, such as buzzers
	EV_REP       uint16 = 0x0014 // Enables autorepeat of keys in the input core
	EV_FF        uint16 = 0x0015 // Sends force-feedback effects to a device
	EV_PWR       uint16 = 0x0016 // Power management events
	EV_FF_STATUS uint16 = 0x0017 // Device reporting of force-feedback effects back to the host
	EV_MAX       uint16 = 0x001F
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
	this.devices = make([]hw.InputDevice,0)
	if err := evFind(func (device *InputDevice) {
		this.devices = append(this.devices,device)
	}); err != nil {
		return nil, err
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

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

// Strinfigy InputDriver object
func (this *InputDriver) String() string {
	return fmt.Sprintf("<linux.Input>{ devices=%v }",this.devices)
}


// Strinfigy InputDevice object
func (this *InputDevice) String() string {
	return fmt.Sprintf("<linux.InputDevice>{ name=\"%s\" path=%s type=%v fd=%v }", this.Name, this.Path, this.Type, this.handle)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Find all input devices
func evFind(callback func (driver *InputDevice)) error {
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

