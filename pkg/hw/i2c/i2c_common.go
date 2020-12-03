package i2c

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type i2c struct {
	gopi.Unit
	sync.Mutex
	gopi.Logger

	devices map[gopi.I2CBus]*device
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *i2c) New(gopi.Config) error {
	this.devices = make(map[gopi.I2CBus]*device, 10)
	return nil
}

func (this *i2c) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var result error
	for bus := range this.devices {
		if err := this.Close(bus); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Get current slave address
func (this *i2c) GetSlave(bus gopi.I2CBus) uint8 {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, exists := this.devices[bus]; exists {
		return device.slave
	} else {
		return 0
	}
}

func (this *i2c) Read(bus gopi.I2CBus) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var buf []byte
	if device, err := this.Open(bus); err != nil {
		return nil, err
	} else if _, err := device.fh.Read(buf); err != nil {
		return nil, err
	} else {
		return buf, nil
	}
}

func (this *i2c) Write(bus gopi.I2CBus, buf []byte) (int, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return 0, err
	} else if n, err := device.fh.Write(buf); err != nil {
		return n, err
	} else {
		return n, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *i2c) Close(bus gopi.I2CBus) error {
	var result error

	device, exists := this.devices[bus]
	if exists {
		this.Debug("i2C Close=>", bus)
		result = device.fh.Close()
		delete(this.devices, bus)
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *i2c) String() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	str := "<i2c"
	if d := this.Devices(); len(d) > 0 {
		str += " bus=" + fmt.Sprint(d)
	}
	for bus, device := range this.devices {
		str += fmt.Sprintf(" device[%v]=%v", bus, device)
	}
	return str + ">"
}
