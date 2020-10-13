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

	for k, v := range this.devices {
		if err := this.Close(v); err != nil {
			result = multierror.Append(result, err)
		}
		this.RWMutex.Lock()
		delete(this.devices, k)
		this.RWMutex.Unlock()
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this Device) equals(other Device) bool {
	return this.Bus == other.Bus && this.Slave == other.Slave
}

func (this *Devices) get(bus, slave uint) gopi.SPI {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Iterate through map to get an open device
	k := Device{bus, slave}
	for other, v := range this.devices {
		if other.equals(k) {
			return v
		}
	}

	// Not found
	return nil
}

func (this *Devices) delete(bus, slave uint) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Delete a key from the map
	delete(this.devices, Device{bus, slave})
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Devices) String() string {
	str := "<spi"
	for _, device := range this.Enumerate() {
		str += " " + fmt.Sprint(device)
	}
	return str + ">"
}
