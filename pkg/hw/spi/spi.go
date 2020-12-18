package spi

import (
	// Frameworks
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type spi struct {
	gopi.Unit
	sync.Mutex

	devices map[gopi.SPIBus]*device
}

////////////////////////////////////////////////////////////////////////////////
// Globals

const (
	maxBus = 9 // Maximum bus number
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *spi) New(gopi.Config) error {
	this.devices = make(map[gopi.SPIBus]*device, maxBus)
	return nil
}

func (this *spi) Dispose() error {
	// Close devices
	var result error
	for bus := range this.devices {
		if err := this.Close(bus); err != nil {
			result = multierror.Append(result, err)
		}
	}

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Release devices
	this.devices = nil

	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *spi) String() string {
	str := "<spi"
	for bus, device := range this.devices {
		str += fmt.Sprintf(" device[%v]=%v", bus, device)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *spi) Open(bus gopi.SPIBus) (*device, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if d, exists := this.devices[bus]; exists {
		return d, nil
	}
	if d, err := NewDevice(bus, 0); err != nil {
		return nil, err
	} else {
		this.devices[bus] = d
		return d, nil
	}
}

func (this *spi) Close(bus gopi.SPIBus) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if d, exists := this.devices[bus]; exists == false {
		return nil
	} else {
		delete(this.devices, bus)
		return d.Close()
	}
}
