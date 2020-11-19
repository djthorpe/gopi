// +build linux

package i2c

import (
	"os"

	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2CFunction linux.I2CFunction

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	// Minimum and maximum bus numbers
	minBus, maxBus = 0, 9
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Return all valid devices
func (this *i2c) Devices() []gopi.I2CBus {
	var devices []gopi.I2CBus
	for bus := uint(minBus); bus <= uint(maxBus); bus++ {
		if _, err := os.Stat(linux.I2CDevice(bus)); err == nil {
			devices = append(devices, gopi.I2CBus(bus))
		}
	}
	return devices
}

// Set current slave address
func (this *i2c) SetSlave(bus gopi.I2CBus, slave uint8) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return err
	} else if err := linux.I2CSetSlave(device.Fd(), slave); err != nil {
		return err
	} else {
		device.slave = slave
	}

	// Return success
	return nil
}

// Return true if a slave was detected at a particular address
func (this *i2c) DetectSlave(bus gopi.I2CBus, slave uint8) (bool, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return false, err
	} else if result, err := linux.I2CDetectSlave(device.Fd(), slave, linux.I2CFunction(device.funcs)); err != nil {
		return false, err
	} else {
		return result, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - READ

func (this *i2c) ReadUint8(bus gopi.I2CBus, reg uint8) (uint8, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return 0, err
	} else if result, err := linux.I2CReadUint8(device.Fd(), reg, linux.I2CFunction(device.funcs)); err != nil {
		return 0, err
	} else {
		return result, nil
	}
}

func (this *i2c) ReadInt8(bus gopi.I2CBus, reg uint8) (int8, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return 0, err
	} else if result, err := linux.I2CReadInt8(device.Fd(), reg, linux.I2CFunction(device.funcs)); err != nil {
		return 0, err
	} else {
		return result, nil
	}
}

func (this *i2c) ReadUint16(bus gopi.I2CBus, reg uint8) (uint16, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return 0, err
	} else if result, err := linux.I2CReadUint16(device.Fd(), reg, linux.I2CFunction(device.funcs)); err != nil {
		return 0, err
	} else {
		return result, nil
	}
}

func (this *i2c) ReadInt16(bus gopi.I2CBus, reg uint8) (int16, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return 0, err
	} else if result, err := linux.I2CReadInt16(device.Fd(), reg, linux.I2CFunction(device.funcs)); err != nil {
		return 0, err
	} else {
		return result, nil
	}
}

func (this *i2c) ReadBlock(bus gopi.I2CBus, reg, length uint8) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return nil, err
	} else if result, err := linux.I2CReadBlock(device.Fd(), reg, length, linux.I2CFunction(device.funcs)); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - WRITE

func (this *i2c) WriteUint8(bus gopi.I2CBus, reg, value uint8) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return linux.I2CWriteUint8(device.Fd(), reg, value, linux.I2CFunction(device.funcs))
	}
}

func (this *i2c) WriteInt8(bus gopi.I2CBus, reg uint8, value int8) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return linux.I2CWriteInt8(device.Fd(), reg, value, linux.I2CFunction(device.funcs))
	}
}

func (this *i2c) WriteUint16(bus gopi.I2CBus, reg uint8, value uint16) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return linux.I2CWriteUint16(device.Fd(), reg, value, linux.I2CFunction(device.funcs))
	}
}

func (this *i2c) WriteInt16(bus gopi.I2CBus, reg uint8, value int16) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return linux.I2CWriteInt16(device.Fd(), reg, value, linux.I2CFunction(device.funcs))
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f I2CFunction) String() string {
	return linux.I2CFunction(f).String()
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *i2c) Open(bus gopi.I2CBus) (*device, error) {
	device, exists := this.devices[bus]
	if exists == false {
		if fh, err := linux.I2COpenDevice(uint(bus)); err != nil {
			return nil, err
		} else if funcs, err := linux.I2CFunctions(fh.Fd()); err != nil {
			return nil, err
		} else {
			device = NewDevice(bus, fh, I2CFunction(funcs))
			this.devices[bus] = device
		}
	}
	return device, nil
}
