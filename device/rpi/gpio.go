/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

import (
	"bytes"
	"encoding/binary"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

type GPIO struct {
	base    uint32
	memlock sync.Mutex
	mem8    []uint8
	mem     []uint32
}

type Pin uint8
type State uint8
type Direction uint8
type Pull uint8

////////////////////////////////////////////////////////////////////////////////

const (
	GPIO_DEV_GPIOMEM        = "/dev/gpiomem"
	GPIO_DEV_MEM            = "/dev/mem"
	GPIO_BASE        uint32 = 0x200000
	GPIO_MEMLENGTH          = 4096
	PINMASK          uint32 = 7 // pin mode is 3 bits
)

// Pin direction
const (
	INPUT Direction = iota
	OUTPUT
)

// Pin state
const (
	LOW State = iota
	HIGH
)

// Pull Up / Down / Off
const (
	PULLOFF Pull = iota
	PULLDOWN
	PULLUP
)

////////////////////////////////////////////////////////////////////////////////

// Create new GPIO object, returns error if not possible
func (rpi *RaspberryPi) NewGPIO() (*GPIO, error) {
	var file *os.File
	var err error
	var base uint32

	// open GPIO file
	if file, err = os.OpenFile(GPIO_DEV_GPIOMEM, os.O_RDWR|os.O_SYNC, 0); os.IsNotExist(err) {
		file, err = os.OpenFile(GPIO_DEV_MEM, os.O_RDWR|os.O_SYNC, 0)
		if err != nil {
			return nil, err
		}
		base, err = getBaseAddress(rpi)
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create object, and lock object
	this := new(GPIO)
	this.memlock.Lock()
	defer this.memlock.Unlock()

	// Memory map GPIO registers to byte array
	this.mem8, err = syscall.Mmap(int(file.Fd()), int64(base), GPIO_MEMLENGTH, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&this.mem8))
	header.Len /= (32 / 8) // (32 bit = 4 bytes)
	header.Cap /= (32 / 8)
	this.mem = *(*[]uint32)(unsafe.Pointer(&header))

	return this, nil
}

// Close unmaps GPIO memory
func (this *GPIO) Close() error {
	this.memlock.Lock()
	defer this.memlock.Unlock()
	return syscall.Munmap(this.mem8)
}

////////////////////////////////////////////////////////////////////////////////

// WritePin sets a given pin HIGH or LOW
func (this *GPIO) WritePin(pin Pin, state State) {

	p := uint8(pin)

	// Clear register, 10 / 11 depending on bank
	// Set register, 7 / 8 depending on bank
	clearReg := p/32 + 10
	setReg := p/32 + 7

	this.memlock.Lock()
	defer this.memlock.Unlock()

	if state == LOW {
		this.mem[clearReg] = 1 << (p & 31)
	} else {
		this.mem[setReg] = 1 << (p & 31)
	}
}

// WritePinLow sets a pin to LOW
func (this *GPIO) WritePinLow(pin Pin) {
	this.WritePin(pin, LOW)
}

// WritePinHigh sets a pin to HIGH
func (this *GPIO) WritePinHigh(pin Pin) {
	this.WritePin(pin, HIGH)
}

// ReadPin reads the state of a pin
func (this *GPIO) ReadPin(pin Pin) State {
	// Input level register offset (13 / 14 depending on bank)
	levelReg := uint8(pin)/32 + 13

	if (this.mem[levelReg] & (1 << uint8(pin))) != 0 {
		return HIGH
	}
	return LOW
}

// SetPinMode sets the direction of a given pin (INPUT or OUTPUT)
func (this *GPIO) SetPinMode(pin Pin, direction Direction) {

	// Pin fsel register, 0 or 1 depending on bank
	fsel := uint8(pin) / 10
	shift := (uint8(pin) % 10) * 3

	this.memlock.Lock()
	defer this.memlock.Unlock()

	if direction == INPUT {
		this.mem[fsel] = this.mem[fsel] &^ (PINMASK << shift)
	} else {
		this.mem[fsel] = (this.mem[fsel] &^ (PINMASK << shift)) | (1 << shift)
	}
}

// SetPinModeInput sets pin to INPUT
func (this *GPIO) SetPinModeInput(pin Pin) {
	this.SetPinMode(pin, INPUT)
}

// SetPinModeOutput sets pin to OUTPUT
func (this *GPIO) SetPinModeOutput(pin Pin) {
	this.SetPinMode(pin, OUTPUT)
}

////////////////////////////////////////////////////////////////////////////////

// Read /proc/device-tree/soc/ranges and determine the base address.
// Use the default Raspberry Pi 1 base address if this fails.
func getBaseAddress(pi *RaspberryPi) (uint32, error) {
	peripheralbase, err := pi.PeripheralBase()
	if err != nil {
		return 0, err
	}
	ranges, err := os.Open("/proc/device-tree/soc/ranges")
	if err != nil {
		return uint32(peripheralbase + GPIO_BASE), nil
	}
	defer ranges.Close()
	b := make([]byte, 4)
	n, err := ranges.ReadAt(b, 4)
	if n != 4 || err != nil {
		return uint32(peripheralbase + GPIO_BASE), nil
	}
	buf := bytes.NewReader(b)
	var out uint32
	err = binary.Read(buf, binary.BigEndian, &out)
	if err != nil {
		return uint32(peripheralbase + GPIO_BASE), nil
	}
	return uint32(out + GPIO_BASE), nil
}
