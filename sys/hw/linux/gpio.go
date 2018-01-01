// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GPIO struct {
	UnexportOnClose bool
	FilePoll        FilePollInterface
}

type gpio struct {
	log         gopi.Logger
	exported    []gopi.GPIOPin
	watched     map[gopi.GPIOPin]*os.File
	lock        sync.Mutex
	filepoll    FilePollInterface
	subscribers []chan gopi.Event
}

type event struct {
	driver *gpio
	pin    gopi.GPIOPin
	edge   gopi.GPIOEdge
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
	logger.Debug("sys.hw.linux.GPIO.Open{ UnexportOnClose=%v }", config.UnexportOnClose)

	this := new(gpio)
	this.log = logger
	this.watched = make(map[gopi.GPIOPin]*os.File, 0)

	// Make array of exported pins
	if config.UnexportOnClose {
		this.exported = make([]gopi.GPIOPin, 0)
	}

	// File Poll module is required or else returns ErrBadParameter
	if config.FilePoll != nil {
		this.filepoll = config.FilePoll
	} else {
		return nil, gopi.ErrBadParameter
	}

	// Subscribers
	this.subscribers = make([]chan gopi.Event, 0)

	// Success
	return this, nil
}

// Close
func (this *gpio) Close() error {
	this.log.Debug("sys.hw.linux.GPIO.Close{ }")

	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()

	// Unwatch pins
	for pin, file := range this.watched {
		this.filepoll.Unwatch(file)
		if err := file.Close(); err != nil {
			this.log.Warn("sys.hw.linux.GPIO.Close: %v: %v", pin, err)
		}
	}

	if len(this.exported) > 0 {
		// unexport pins
		for _, pin := range this.exported {
			if isExported(pin) {
				if err := unexportPin(pin); err != nil {
					this.log.Warn("sys.hw.linux.GPIO.Close: Unable to export pin %v: %v", pin, err)
				}
			}
		}
	}

	// Close subscriber channels
	for _, c := range this.subscribers {
		if c != nil {
			close(c)
		}
	}

	// Zero out member variables
	this.exported = nil
	this.watched = nil
	this.filepoll = nil
	this.subscribers = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - RETURN INFORMATION

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

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - READ/WRITE

// Read pin state
func (this *gpio) ReadPin(pin gopi.GPIOPin) gopi.GPIOState {
	this.log.Debug2("<sys.hw.linux.GPIO.ReadPin>{ pin=%v }", pin)

	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()

	// Check for pin exported
	if err := this.exportPin(pin); err != nil {
		this.log.Error("Unable to export %v: %v", pin, err)
		return gopi.GPIO_LOW
	}
	// Do extra checks of output state when debugging is on
	if this.log.IsDebug() {
		if direction, err := direction(pin); err != nil {
			this.log.Warn("Invalid direction for pin %v: '%v'", pin, err)
		} else if direction != "in" {
			this.log.Warn("Invalid direction for pin %v: '%v'", pin, direction)
		}
	}
	// Read the pin
	if value, err := readPin(pin); err != nil {
		this.log.Error("Unable to read %v: %v", pin, err)
		return gopi.GPIO_LOW
	} else {
		switch value {
		case "0":
			return gopi.GPIO_LOW
		case "1":
			return gopi.GPIO_HIGH
		default:
			this.log.Warn("Invalid value for pin %v: '%v'", pin, value)
			return gopi.GPIO_HIGH
		}
	}
}

// WritePin writes pin state - either low or high
func (this *gpio) WritePin(pin gopi.GPIOPin, state gopi.GPIOState) {
	this.log.Debug2("<sys.hw.linux.GPIO.WritePin>{ pin=%v state=%v }", pin, state)

	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()

	// Check for pin exported
	if err := this.exportPin(pin); err != nil {
		this.log.Error("Unable to export %v: %v", pin, err)
		return
	}
	// Do extra checks of output state when debugging is on
	if this.log.IsDebug() {
		if direction, err := direction(pin); err != nil {
			this.log.Warn("Invalid pin direction for %v: '%v'", pin, err)
		} else if direction != "out" {
			this.log.Warn("Invalid pin direction for %v: '%v'", pin, direction)
		}
	}
	// Write pin
	switch state {
	case gopi.GPIO_LOW:
		if err := writePin(pin, "0"); err != nil {
			this.log.Error("Unable to write value to %v: %v", pin, err)
		}
	case gopi.GPIO_HIGH:
		if err := writePin(pin, "1"); err != nil {
			this.log.Error("Unable to write value to %v: %v", pin, err)
		}
	}
}

// GetPinMode gets pin mode, which is either in or out
// or returns GPIO_NONE on error
func (this *gpio) GetPinMode(pin gopi.GPIOPin) gopi.GPIOMode {
	this.log.Debug2("<sys.hw.linux.GPIO.GetPinMode>{ pin=%v }", pin)

	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()

	// Check for pin exported
	if err := this.exportPin(pin); err != nil {
		this.log.Error("Unable to export %v: %v", pin, err)
		return gopi.GPIO_NONE
	}
	// Read the pin
	if value, err := direction(pin); err != nil {
		this.log.Error("Unable to read direction %v: %v", pin, err)
		return gopi.GPIO_NONE
	} else {
		switch value {
		case "in":
			return gopi.GPIO_INPUT
		case "out":
			return gopi.GPIO_OUTPUT
		default:
			this.log.Warn("Invalid direction for %v: '%v'", pin, value)
			return gopi.GPIO_NONE
		}
	}
}

// SetPinMode set pin mode to either in or out. No other
// modes are supported through this driver
func (this *gpio) SetPinMode(pin gopi.GPIOPin, mode gopi.GPIOMode) {
	this.log.Debug2("<sys.hw.linux.GPIO.SetPinMode>{ pin=%v mode=%v }", pin, mode)

	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()

	// Check for pin exported
	if err := this.exportPin(pin); err != nil {
		this.log.Error("Unable to export %v: %v", pin, err)
		return
	}
	// Write pin
	switch mode {
	case gopi.GPIO_INPUT:
		if err := setDirection(pin, "in"); err != nil {
			this.log.Error("Unable to write direction to %v: %v", pin, err)
		}
		if err := writeEdge(pin, "none"); err != nil {
			this.log.Error("Unable to write edge to %v: %v", pin, err)
		}
	case gopi.GPIO_OUTPUT:
		if err := setDirection(pin, "out"); err != nil {
			this.log.Error("Unable to write direction to %v: %v", pin, err)
		}
	default:
		this.log.Error("Invalid pin mode %v: %v", pin, mode)
	}
}

// SetPullMode is not implemented in the sysfs driver
func (this *gpio) SetPullMode(pin gopi.GPIOPin, pull gopi.GPIOPull) error {
	return gopi.ErrNotImplemented
}

// Watch will watch a pin for rising, falling or both edges. When
// set to EDGE_NONE then watching is stopped
func (this *gpio) Watch(pin gopi.GPIOPin, edge gopi.GPIOEdge) error {
	this.log.Debug2("<sys.hw.linux.GPIO.Watch>{ pin=%v edge=%v }", pin, edge)

	// Mutex
	this.lock.Lock()
	defer this.lock.Unlock()

	// Check for pin exported
	if err := this.exportPin(pin); err != nil {
		this.log.Error("Watch: unable to export %v: %v", pin, err)
		return err
	}

	// Do extra checks of output state when debugging is on
	if this.log.IsDebug() {
		if direction, err := direction(pin); err != nil {
			this.log.Warn("Watch: Invalid direction for %v: '%v'", pin, err)
		} else if direction != "in" {
			this.log.Warn("Watch: Invalid direction for %v: '%v'", pin, direction)
		}
	}

	// Set rising, falling, both or none
	edge_write := ""
	switch edge {
	case gopi.GPIO_EDGE_NONE:
		if err := writeEdge(pin, "none"); err != nil {
			this.log.Error("Watch: Unable to write edge for %v: %v", pin, err)
		} else if file, exists := this.watched[pin]; exists == false {
			// IGNORE UNWATCHED PINS
		} else if err := this.filepoll.Unwatch(file); err != nil {
			this.log.Error("%v: %v", pin, err)
			file.Close()
		} else if err := file.Close(); err != nil {
			return err
		} else {
			// Remove from list of watched pins
			delete(this.watched, pin)
		}
	case gopi.GPIO_EDGE_RISING:
		edge_write = "rising"
	case gopi.GPIO_EDGE_FALLING:
		edge_write = "falling"
	case gopi.GPIO_EDGE_BOTH:
		edge_write = "both"
	default:
		return errors.New("Watch: Invalid edge value")
	}

	if edge_write != "" {
		if err := writeEdge(pin, edge_write); err != nil {
			this.log.Error("Watch: Unable to write edge for %v: %v", pin, err)
			return err
		} else if _, exists := this.watched[pin]; exists {
			// IGNORE ALREADY WATCHED PINS
		} else if file, err := watchValue(pin); err != nil {
			this.log.Error("Watch: Unable to watch %v: %v", pin, err)
			return err
		} else if err := this.filepoll.Watch(file, FILEPOLL_MODE_EDGE, func(handle *os.File, mode FilePollMode) {
			if err := this.handleEdge(handle, pin); err != nil {
				this.log.Warn("Watch: %v: %v", pin, err)
			}
		}); err != nil {
			this.log.Error("Watch: %v: %v", pin, err)
			file.Close()
			return err
		} else {
			this.watched[pin] = file
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - EVENTS

// Subscribe to events emitted. Returns unique subscriber
// identifier and channel on which events are emitted
func (this *gpio) Subscribe() chan gopi.Event {
	this.log.Debug2("<sys.hw.linux.GPIO.Subscribe>{ }")

	// Create a new channel for emitting events
	subscriber := make(chan gopi.Event)
	this.subscribers = append(this.subscribers, subscriber)
	return subscriber
}

// Unsubscribe from events emitted
func (this *gpio) Unsubscribe(subscriber chan gopi.Event) {
	this.log.Debug2("<sys.hw.linux.GPIO.Unsubscribe>{ }")

	for i := range this.subscribers {
		if this.subscribers[i] == subscriber {
			this.subscribers[i] = nil
			close(subscriber)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE - EVENT

func (this *event) Name() string {
	return "GPIOEvent"
}

func (this *event) Source() gopi.Driver {
	return this.driver
}

func (this *event) Pin() gopi.GPIOPin {
	return this.pin
}

func (this *event) Edge() gopi.GPIOEdge {
	return this.edge
}

func (this *event) String() string {
	return fmt.Sprintf("<sys.hw.linux.GPIO.Event>{ pin=%v edge=%v source=%v }", this.pin, this.edge, this.driver)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *gpio) exportPin(pin gopi.GPIOPin) error {
	if isExported(pin) == false {
		if err := exportPin(pin); err != nil {
			return err
		}
	}
	if this.exported != nil {
		for _, exported := range this.exported {
			if pin == exported {
				return nil
			}
		}
		this.exported = append(this.exported, pin)
	}
	return nil
}

func (this *gpio) handleEdge(handle *os.File, pin gopi.GPIOPin) error {
	if _, err := handle.Seek(0, io.SeekStart); err != nil {
		return err
	} else if buf, err := ioutil.ReadAll(handle); err != nil {
		return err
	} else {
		value := strings.TrimSpace(string(buf))
		switch value {
		case "0":
			this.emit(pin, gopi.GPIO_EDGE_FALLING)
		case "1":
			this.emit(pin, gopi.GPIO_EDGE_RISING)
		default:
			this.emit(pin, gopi.GPIO_EDGE_NONE)
		}
		return nil
	}
}

// Emit an event
func (this *gpio) emit(pin gopi.GPIOPin, edge gopi.GPIOEdge) {
	event := &event{driver: this, pin: pin, edge: edge}
	for _, channel := range this.subscribers {
		if channel != nil {
			channel <- event
		}
	}
}

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
	if err := writeFile(GPIO_EXPORT, strconv.FormatUint(uint64(pin), 10)+"\n"); err != nil {
		return err
	} else {
		// Wait for 50ms for things to settle
		time.Sleep(50 * time.Millisecond)
	}
	// check to make sure pin is exported
	if isExported(pin) == false {
		return fmt.Errorf("exportPin %v failed", pin)
	}
	// Set edge to 'none' if direction is 'in'
	if dir, err := direction(pin); err != nil {
		return err
	} else if dir == "in" {
		if err := writeEdge(pin, "none"); err != nil {
			return err
		}
	}
	// Success
	return nil
}

func unexportPin(pin gopi.GPIOPin) error {
	return writeFile(GPIO_UNEXPORT, strconv.FormatUint(uint64(pin), 10)+"\n")
}

func direction(pin gopi.GPIOPin) (string, error) {
	if value, err := readFile(filenameForPin(pin, "direction")); err != nil {
		return "", err
	} else {
		return strings.TrimSpace(value), nil
	}
}

func setDirection(pin gopi.GPIOPin, value string) error {
	return writeFile(filenameForPin(pin, "direction"), value+"\n")
}

func readPin(pin gopi.GPIOPin) (string, error) {
	if value, err := readFile(filenameForPin(pin, "value")); err != nil {
		return "", err
	} else {
		return strings.TrimSpace(value), nil
	}
}

func writePin(pin gopi.GPIOPin, value string) error {
	return writeFile(filenameForPin(pin, "value"), value+"\n")
}

func writeEdge(pin gopi.GPIOPin, edge string) error {
	return writeFile(filenameForPin(pin, "edge"), edge+"\n")
}

func readEdge(pin gopi.GPIOPin) (string, error) {
	if value, err := readFile(filenameForPin(pin, "edge")); err != nil {
		return "", err
	} else {
		return strings.TrimSpace(value), nil
	}
}

func watchValue(pin gopi.GPIOPin) (*os.File, error) {
	return os.OpenFile(filenameForPin(pin, "value"), os.O_RDONLY, 0)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *gpio) String() string {
	str := ""
	if this.exported != nil {
		str = str + fmt.Sprintf(" exportedPins=%v", this.exported)
	}
	if len(this.watched) > 0 {
		str = str + " watchedPins=["
		for k := range this.watched {
			str = str + " " + fmt.Sprint(k)
		}
		str = str + " ]"
	}
	return fmt.Sprintf("sys.hw.linux.GPIO{%v }", str)
}
