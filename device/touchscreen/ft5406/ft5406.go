/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This package implements the concrete interface to the FT5406 device
// for for Raspberry Pi. The FT5406 is the touchscreen which comes as
// part of the official Raspberry Pi LCD screen
package ft5406

// System imports
import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// Local imports
import (
	"../../../input" /* Absolute: github.com/djthorpe/gopi/input */
)

////////////////////////////////////////////////////////////////////////////////

// The configuration options for the FT5406. Basically, there are no
// configuration options for this device.
type Config struct { }

// The driver state
type Driver struct {
	device string
	name   string
	file   *os.File
}

////////////////////////////////////////////////////////////////////////////////

const (
	PATH_INPUT_DEVICES   = "/sys/class/input/event*"
	MAX_SLOTS uint       = 10
)

////////////////////////////////////////////////////////////////////////////////

var (
	REGEXP_DEVICENAME = regexp.MustCompile("^FT5406")
)

////////////////////////////////////////////////////////////////////////////////
// input.Opener interface

// Concrete Open method
func (config Config) Open() (input.Driver, error) {
	return NewFT5406()
}

////////////////////////////////////////////////////////////////////////////////
// Create Touchscreen

func NewFT5406() (*Driver, error) {
	var err error

	this := new(Driver)
	this.name, this.device, err = getDeviceNameAndPath()
	if err != nil {
		return nil, err
	}

	// open driver
	this.file, err = os.Open(this.device)
	if err != nil {
		return nil, err
	}

	return this, nil
}

func (this *Driver) Close() error {
	return this.file.Close()
}

func (this *Driver) GetName() string {
	return this.name
}

func (this *Driver) GetType() input.DeviceType {
	return input.TYPE_TOUCHSCREEN
}

func (this *Driver) GetFd() *os.File {
	return this.file
}

func (this *Driver) GetSlots() uint {
	return MAX_SLOTS
}

////////////////////////////////////////////////////////////////////////////////
// Private Methods

func getDeviceNameAndPath() (string, string, error) {
	files, err := filepath.Glob(PATH_INPUT_DEVICES)
	if err != nil {
		return "", "", err
	}
	for _, file := range files {
		buf, err := ioutil.ReadFile(path.Join(file, "device", "name"))
		if err != nil {
			continue
		}
		if REGEXP_DEVICENAME.Match(buf) {
			return strings.TrimSpace(string(buf)), path.Join("/", "dev", "input", path.Base(file)), nil
		}
	}
	return "", "", errors.New("Device not found")
}
