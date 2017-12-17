/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GPIO struct{}

type gpio struct {
	log gopi.Logger
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GPIO_EXPORT   = "/sys/class/gpio/export"
	GPIO_UNEXPORT = "/sys/class/gpio/unexport"
	GPIO_PIN      = "/sys/class/gpio/gpio%v"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config GPIO) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.hw.linux.GPIO.Open{ }")

	this := new(gpio)
	this.log = logger

	// Success
	return this, nil
}

// Close
func (this *gpio) Close() error {
	this.log.Debug("sys.hw.linux.GPIO.Close{ }")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

// Return number of physical pins, or 0 if if cannot be returned
func (this *gpio) NumberOfPhysicalPins() uint {
	return 0
}

// Return array of available logical pins
func (this *gpio) Pins() []gopi.GPIOPin {
	return []gopi.GPIOPin{}
}

// Return logical pin for physical pin number. Returns
// GPIO_PIN_NONE where there is no logical pin at that position
func (this *gpio) PhysicalPin(uint) gopi.GPIOPin {
	return gopi.GPIO_PIN_NONE
}

// Return physical pin number for logical pin. Returns 0 where there
// is no physical pin for this logical pin
func (this *gpio) PhysicalPinForPin(gopi.GPIOPin) uint {
	return 0
}

// Read pin state
func (this *gpio) ReadPin(gopi.GPIOPin) gopi.GPIOState {
	return gopi.GPIO_LOW
}

// WritePin Writes pin state
func (this *gpio) WritePin(pin gopi.GPIOPin, state gopi.GPIOState) {
	// Check for pin exported
	if isExported(pin) == false {
		if err := exportPin(pin); err != nil {
			this.log.Error("Unable to export %v: %v", pin, err)
			return
		}
	}
	// Write pin
	switch state {
	case gopi.GPIO_LOW:
		if err := writePin(pin, "0"); err != nil {
			this.log.Error("Unable to write to pin %v: %v", pin, err)
		}
	case gopi.GPIO_HIGH:
		if err := writePin(pin, "1"); err != nil {
			this.log.Error("Unable to write to pin %v: %v", pin, err)
		}
	}
}

// Get pin mode
func (this *gpio) GetPinMode(gopi.GPIOPin) gopi.GPIOMode {
	return gopi.GPIO_ALT0
}

// Set pin mode
func (this *gpio) SetPinMode(gopi.GPIOPin, gopi.GPIOMode) {
	// TODO
}

// Set pull mode
func (this *gpio) SetPullMode(gopi.GPIOPin, gopi.GPIOPull) {
	// TODO
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func writeFile(filename string, value string) error {
	return ioutil.WriteFile(filename, []byte(value), 777)
}

func readFile(filename string) (string, error) {
	if bytes, err := ioutil.ReadFile(filename); err != nil {
		return "", err
	} else {
		return string(bytes), nil
	}
}

func isExported(pin gopi.GPIOPin) bool {
	if _, err := os.Stat(filenameForPin(pin, "")); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	} else {
		return true
	}
}

func filenameForPin(pin gopi.GPIOPin, filename string) string {
	return filepath.Join(fmt.Sprintf(GPIO_PIN, uint(pin)), filename)
}

func exportPin(pin gopi.GPIOPin) error {
	return writeFile(GPIO_EXPORT, strconv.FormatUint(uint64(pin), 10))
}

func unexportPin(pin gopi.GPIOPin) error {
	return writeFile(GPIO_UNEXPORT, strconv.FormatUint(uint64(pin), 10))
}

func direction(pin gopi.GPIOPin) (string, error) {
	return readFile(filenameForPin(pin, "direction"))
}

func setDirection(pin gopi.GPIOPin, direction string) error {
	return writeFile(filenameForPin(pin, "direction"), direction)
}

func readPin(pin gopi.GPIOPin) (string, error) {
	return readFile(filenameForPin(pin, "value"))
}

func writePin(pin gopi.GPIOPin, value string) error {
	return writeFile(filenameForPin(pin, "value"), value)
}
