package spi

import (
	// Frameworks
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Devices struct {
	gopi.Unit
	sync.RWMutex

	devices map[Device]gopi.SPI
}

type Device struct {
	Bus, Slave uint
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// New is called to initialize
func (this *Devices) New(gopi.Config) error {
	this.devices = make(map[Device]gopi.SPI)
	return nil
}

// Dispose is called to close
func (this *Devices) Dispose() error {
	var result error

	for k := range this.devices {
		if err := this.Close(k); err != nil {
			result = multierror.Append(result, err)
		}
		this.RWMutex.Lock()
		delete(this.devices, k)
		this.RWMutex.Unlock()
	}

	return result
}

// Get returns a gopi.SPI object based on bus and slave
func (this *Devices) Get(device Device) gopi.SPI {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Iterate to get device
	for k, v := range this.devices {
		if device.equals(k) {
			return v
		}
	}

	// Not found
	return nil
}

func (this *Devices) Delete(device Device) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	delete(this.devices, device)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this Device) equals(other Device) bool {
	return this.Bus == other.Bus && this.Slave == other.Slave
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Devices) String() string {
	str := "<spi"
	for _, device := range this.Enumerate() {
		str += " " + fmt.Sprint(device) + "=" + fmt.Sprint(this.Get(device))
	}
	return str + ">"
}

func (this Device) String() string {
	str := "<device"
	str += " bus=" + fmt.Sprint(this.Bus)
	str += " slave=" + fmt.Sprint(this.Slave)
	return str + ">"
}
