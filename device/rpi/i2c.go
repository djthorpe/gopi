/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
	"sync"
	"syscall"
	"reflect"
	"unsafe"
	"os"
	"errors"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2C struct {
	Device gopi.HardwareDriver
	Master uint
}

type I2CDriver struct {
	log      *util.LoggerDevice // logger
	memlock  sync.Mutex
	master   uint
	mem8     []uint8             // access I2C registers as bytes
	mem32    []uint32            // access I2C registers as uint32
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RPI_I2C_SDA1 = 3 // Physical pin 3
	RPI_I2C_SCL1 = 5 // Physical pin 5
)

const (
	I2C_DEV_I2CMEM        = "/dev/i2c"
	I2C_DEV_MEM           = "/dev/mem"
	I2C_BASE_M0       uint32 = 0x00205000
	I2C_BASE_M1       uint32 = 0x00804000
	I2C_BASE_M2       uint32 = 0x00805000
	I2C_SIZE          uint32 = 32       // 8 registers (8 * 4 bytes)
	I2C_REG_CTRL      uint32 = 0x00     // Control
	I2C_REG_STATUS    uint32 = 0x04     // Status
	I2C_REG_DLEN      uint32 = 0x08     // Data Length
	I2C_REG_ADDR      uint32 = 0x0C     // Slave Address
	I2C_REG_FIFO      uint32 = 0x10     // Data FIFO
	I2C_REG_DIV       uint32 = 0x14     // Clock Divider
	I2C_REG_DEL       uint32 = 0x18     // Data Delay
	I2C_REG_CLKT      uint32 = 0x1C     // Clock Stretch Timeout
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new I2C object, returns error if not possible
func (config I2C) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	var err error
	log.Debug("<rpi.I2C>Open")

	// create new GPIO driver
	this := new(I2CDriver)

	// Set logging & device
	this.log = log
	this.master = config.Master
	
	// Lock memory
	this.memlock.Lock()
	defer this.memlock.Unlock()

	// Open the /dev/mem and provide offset & size for accessing memory
	file, peripheral_base, peripheral_size, err := i2cOpenDevice(config.Device.(*DeviceState),config.Master)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Memory map
	log.Info("base=%08X size=%08X",peripheral_base,peripheral_size)
	this.mem8, err = syscall.Mmap(int(file.Fd()),int64(peripheral_base), int(peripheral_size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&this.mem8))
	header.Len /= (32 / 8)
	header.Cap /= (32 / 8)
	this.mem32 = *(*[]uint32)(unsafe.Pointer(&header))

	// success
	return this, nil
}

// Close I2C connection
func (this *I2CDriver) Close() error {
	this.log.Debug("<rpi.I2C>Close")

	// Unmap memory and return error
	this.memlock.Lock()
	defer this.memlock.Unlock()
	return syscall.Munmap(this.mem8)
}

// Strinfigy I2C object
func (this *I2CDriver) String() string {
	return fmt.Sprintf("<rpi.I2C>{ master=%v }", this.master)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func i2cOpenDevice(device *DeviceState,master uint) (*os.File, uint32, uint32, error) {
	var file *os.File
	var err error

	// Calculate peripheral_base
	peripheral_base := device.GetPeripheralAddress()
	switch(master) {
		case 0:
			peripheral_base = peripheral_base + I2C_BASE_M0
			break
		case 1:
			peripheral_base = peripheral_base + I2C_BASE_M1
			break
		case 2:
			peripheral_base = peripheral_base + I2C_BASE_M2
			break
		default:
			return nil, 0, 0, errors.New("Invalid I2C master number")
	}

	// open memory
	file, err = os.OpenFile(I2C_DEV_MEM, os.O_RDWR|os.O_SYNC, 0)
	if err != nil {
		return nil, 0, 0, err
	}

	// success
	return file, peripheral_base, I2C_SIZE, nil
}




