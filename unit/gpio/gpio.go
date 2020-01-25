package gpio

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type GPIO struct{}

type gpio struct {
	sysfs gopi.GPIO
	rpi   gopi.GPIO

	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (GPIO) Name() string { return "gopi/gpio" }

func (config GPIO) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(gpio)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// INIT & CLOSE

func (this *gpio) Init(config GPIO) error {

	// Return success
	return nil
}

func (this *gpio) Close() error {

	// Release resources
	this.sysfs = nil
	this.rpi = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *gpio) String() string {
	str := "<" + this.Log.Name()
	if this.sysfs != nil {
		str += " sysfs=" + this.sysfs.String()
	}
	if this.rpi != nil {
		str += " rpi=" + this.rpi.String()
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.GPIO

// Return number of physical pins, or 0 if if cannot be returned
// or nothing is known about physical pins
func (this *gpio) NumberOfPhysicalPins() uint {
	switch {
	case this.rpi != nil:
		return this.rpi.NumberOfPhysicalPins()
	case this.sysfs != nil:
		return this.sysfs.NumberOfPhysicalPins()
	default:
		return 0
	}
}

// Return array of available logical pins or nil if nothing is
// known about pins
func (this *gpio) Pins() []gopi.GPIOPin {
	switch {
	case this.rpi != nil:
		return this.rpi.Pins()
	case this.sysfs != nil:
		return this.sysfs.Pins()
	default:
		return nil
	}
}

// Return logical pin for physical pin number. Returns
// GPIO_PIN_NONE where there is no logical pin at that position
// or we don't now about the physical pins
func (this *gpio) PhysicalPin(physical uint) gopi.GPIOPin {
	switch {
	case this.rpi != nil:
		return this.rpi.PhysicalPin(physical)
	case this.sysfs != nil:
		return this.sysfs.PhysicalPin(physical)
	default:
		return gopi.GPIO_PIN_NONE
	}
}

// Return physical pin number for logical pin. Returns 0 where there
// is no physical pin for this logical pin, or we don't know anything
// about the layout
func (this *gpio) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	switch {
	case this.rpi != nil:
		return this.rpi.PhysicalPinForPin(logical)
	case this.sysfs != nil:
		return this.sysfs.PhysicalPinForPin(logical)
	default:
		return 0
	}
}

// Read pin state
func (this *gpio) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	switch {
	case this.rpi != nil:
		return this.rpi.ReadPin(logical)
	case this.sysfs != nil:
		return this.sysfs.ReadPin(logical)
	default:
		return 0
	}
}

// Write pin state
func (this *gpio) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {
	switch {
	case this.rpi != nil:
		this.rpi.WritePin(logical, state)
	case this.sysfs != nil:
		this.sysfs.WritePin(logical, state)
	}
}

// Get pin mode
func (this *gpio) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	switch {
	case this.rpi != nil:
		return this.rpi.GetPinMode(logical)
	case this.sysfs != nil:
		return this.sysfs.GetPinMode(logical)
	default:
		return 0
	}
}

// Set pin mode
func (this *gpio) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {
	switch {
	case this.rpi != nil:
		this.rpi.SetPinMode(logical, mode)
	case this.sysfs != nil:
		this.sysfs.SetPinMode(logical, mode)
	}
}

// Set pull mode to pull down or pull up - will
// return ErrNotImplemented if not supported
func (this *gpio) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) error {
	switch {
	case this.rpi != nil:
		return this.rpi.SetPullMode(logical, pull)
	case this.sysfs != nil:
		return this.sysfs.SetPullMode(logical, pull)
	default:
		return gopi.ErrNotImplemented
	}
}

// Start watching for rising and/or falling edge,
// or stop watching when GPIO_EDGE_NONE is passed.
// Will return ErrNotImplemented if not supported
func (this *gpio) Watch(logical gopi.GPIOPin, edge gopi.GPIOEdge) error {
	switch {
	case this.sysfs != nil:
		return this.sysfs.Watch(logical, edge)
	case this.rpi != nil:
		return this.rpi.Watch(logical, edge)
	default:
		return gopi.ErrNotImplemented
	}
}
