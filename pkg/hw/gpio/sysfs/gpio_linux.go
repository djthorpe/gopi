// +build linux

package sysfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

type GPIO struct {
	gopi.Unit
	sync.RWMutex
	gopi.Publisher
	gopi.Logger
	*Watcher

	// Flags
	unexportOnDispose *bool

	// State
	exported []gopi.GPIOPin
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GPIO_EXPORT   = "/sys/class/gpio/export"
	GPIO_UNEXPORT = "/sys/class/gpio/unexport"
	GPIO_PIN      = "/sys/class/gpio/gpio%v"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *GPIO) Define(cfg gopi.Config) error {
	this.unexportOnDispose = cfg.FlagBool("gpio.unexport-on-dispose", true, "Unexport exported pins on dispose")
	return nil
}

func (this *GPIO) New(gopi.Config) error {
	// Check for export and unexport paths
	if _, err := os.Stat(GPIO_EXPORT); os.IsNotExist(err) {
		return gopi.ErrNotFound.WithPrefix(GPIO_EXPORT)
	} else if _, err := os.Stat(GPIO_UNEXPORT); os.IsNotExist(err) {
		return gopi.ErrNotFound.WithPrefix(GPIO_EXPORT)
	}

	// Create exported array
	if *this.unexportOnDispose {
		this.exported = make([]gopi.GPIOPin, 0)
	}

	// Return success
	return nil
}

func (this *GPIO) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Unexport pins
	var result error
	for _, pin := range this.exported {
		if isExported(pin) {
			if err := unexportPin(pin); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Release resources
	this.exported = nil

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *GPIO) String() string {
	str := "<gpio.sysfs"
	if this.exported != nil {
		str += " exported=" + fmt.Sprint(this.exported)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *GPIO) NumberOfPhysicalPins() uint {
	return 0
}

func (this *GPIO) Pins() []gopi.GPIOPin {
	return nil
}

func (this *GPIO) PhysicalPin(pin uint) gopi.GPIOPin {
	return gopi.GPIO_PIN_NONE
}

func (this *GPIO) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	return 0
}

func (this *GPIO) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for pin exported
	if err := this.exportPin(logical); err != nil {
		this.Debug("ReadPin: ", err)
		return gopi.GPIO_LOW
	}

	// Read the pin
	value, err := readPin(logical)
	if err != nil {
		this.Debug("ReadPin: ", err)
		return gopi.GPIO_LOW
	}

	// Translate the value
	switch value {
	case "0":
		return gopi.GPIO_LOW
	case "1":
		return gopi.GPIO_HIGH
	default:
		this.Debug("ReadPin: Invalid value")
		return gopi.GPIO_LOW
	}
}

func (this *GPIO) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for pin exported
	if err := this.exportPin(logical); err != nil {
		this.Debug("WritePin: ", err)
		return
	}

	// Do extra checks of output state when debugging is on
	if this.Logger.IsDebug() {
		if direction, err := direction(logical); err != nil {
			this.Debug("WritePin: ", err)
		} else if direction != "out" {
			this.Debug("WritePin: ", err)
		}
	}

	// Write pin
	switch state {
	case gopi.GPIO_LOW:
		if err := writePin(logical, "0"); err != nil {
			this.Debug("WritePin: ", err)
		}
	case gopi.GPIO_HIGH:
		if err := writePin(logical, "1"); err != nil {
			this.Debug("WritePin: ", err)
		}
	}
}

func (this *GPIO) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for pin exported
	if err := this.exportPin(logical); err != nil {
		this.Debug("GetPinMode: ", err)
		return gopi.GPIO_NONE
	}

	// Read the pin
	if value, err := direction(logical); err != nil {
		this.Debug("GetPinMode: ", err)
		return gopi.GPIO_NONE
	} else {
		switch value {
		case "in":
			return gopi.GPIO_INPUT
		case "out":
			return gopi.GPIO_OUTPUT
		default:
			this.Debug("GetPinMode: Unexpected value: ", strconv.Quote(value))
			return gopi.GPIO_NONE
		}
	}
}

// Set pin mode
func (this *GPIO) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for pin exported
	if err := this.exportPin(logical); err != nil {
		this.Debug("SetPinMode: ", err)
		return
	}

	// Write pin
	switch mode {
	case gopi.GPIO_INPUT:
		if err := setDirection(logical, "in"); err != nil {
			this.Debug("SetPinMode: ", err)
		}
		if err := writeEdge(logical, "none"); err != nil {
			this.Debug("SetPinMode: ", err)
		}
	case gopi.GPIO_OUTPUT:
		if err := setDirection(logical, "out"); err != nil {
			this.Debug("SetPinMode: ", err)
		}
	default:
		this.Debug("SetPinMode: Unexpected value: ", mode)
	}
}

func (this *GPIO) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) error {
	return gopi.ErrNotImplemented
}

func (this *GPIO) Watch(logical gopi.GPIOPin, edge gopi.GPIOEdge) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for pin exported
	if err := this.exportPin(logical); err != nil {
		return err
	}

	// Do extra checks of output state when debugging is on
	if this.Logger.IsDebug() {
		if direction, err := direction(logical); err != nil {
			return err
		} else if direction != "in" {
			return gopi.ErrOutOfOrder.WithPrefix("Watch")
		}
	}

	// Set rising, falling, both or none
	value := ""
	switch edge {
	case gopi.GPIO_EDGE_NONE:
		if err := writeEdge(logical, "none"); err != nil {
			return err
		} else if err := this.Watcher.Unwatch(logical); err != nil {
			return err
		}
	case gopi.GPIO_EDGE_RISING:
		value = "rising"
	case gopi.GPIO_EDGE_FALLING:
		value = "falling"
	case gopi.GPIO_EDGE_BOTH:
		value = "both"
	default:
		return gopi.ErrBadParameter.WithPrefix("Watch")
	}

	this.Debug("WriteEdge: ", logical, ": ", value)
	if err := writeEdge(logical, value); err != nil {
		return err
	}

	// TODO: Check where pin already exists

	// Watch the pin
	if file, err := watchValue(logical); err != nil {
		return err
	} else if err := this.Watcher.Watch(file.Fd(), logical, edge); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *GPIO) exportPin(pin gopi.GPIOPin) error {
	if isExported(pin) == false {
		this.Debug("Exporting Pin: ", pin)
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

func filenameForPin(pin gopi.GPIOPin, filename string) string {
	return filepath.Join(fmt.Sprintf(GPIO_PIN, uint(pin)), filename)
}

func writeFile(filename string, value string) error {
	return ioutil.WriteFile(filename, []byte(value), os.ModeDevice|os.ModeCharDevice)
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

func exportPin(pin gopi.GPIOPin) error {
	if err := writeFile(GPIO_EXPORT, strconv.FormatUint(uint64(pin), 10)+"\n"); err != nil {
		return err
	}
	// Check for export success
	if _, err := os.Stat(filenameForPin(pin, "")); os.IsNotExist(err) {
		return err
	}
	// Return success
	return nil
}

func unexportPin(pin gopi.GPIOPin) error {
	if isExported(pin) {
		// Reset state
		if err := writeEdge(pin, "none"); err != nil {
			return err
		}
		if err := setDirection(pin, "in"); err != nil {
			return err
		}
	}
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
	return os.OpenFile(filenameForPin(pin, "value"), os.O_RDWR, 0660)
}
