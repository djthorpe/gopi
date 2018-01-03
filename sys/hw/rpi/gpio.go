// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"syscall"
	"unsafe"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GPIO struct {
	Hardware gopi.Hardware
}

type gpio struct {
	log     gopi.Logger
	product *Product
	pins    map[gopi.GPIOPin]uint // map of logical to physical pins
	memlock sync.Mutex
	mem8    []uint8  // access GPIO as bytes
	mem32   []uint32 // access GPIO as uint32
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GPIO_DEV_GPIOMEM        = "/dev/gpiomem"
	GPIO_DEV_MEM            = "/dev/mem"
	GPIO_BASE        uint32 = 0x200000
	GPIO_SIZE        uint32 = 4096
	GPIO_MAXPINS            = 54 // GPIO0 to GPIO53
)

const (
	// GPIO Registers
	GPIO_GPLVL0    = 0x0034 // Register to read pins GPIO0-GPIO31
	GPIO_GPLVL1    = 0x0038 // Register to read pins GPIO32-GPIO53
	GPIO_GPSET0    = 0x001C // Register to write HIGH to pins GPIO0-GPIO31
	GPIO_GPSET1    = 0x0020 // Register to write HIGH to pins GPIO32-GPIO53
	GPIO_GPCLR0    = 0x0028 // Register to write LOW to pins GPIO0-GPIO31
	GPIO_GPCLR1    = 0x002C // Register to write LOW to pins GPIO32-GPIO53
	GPIO_GPFSEL0   = 0x0000 // Pin modes for GPIO0-GPIO9
	GPIO_GPFSEL1   = 0x0004 // Pin modes for GPIO10-GPIO19
	GPIO_GPFSEL2   = 0x0008 // Pin modes for GPIO20-GPIO29
	GPIO_GPFSEL3   = 0x000C // Pin modes for GPIO30-GPIO39
	GPIO_GPFSEL4   = 0x0010 // Pin modes for GPIO40-GPIO49
	GPIO_GPFSEL5   = 0x0014 // Pin modes for GPIO50-GPIO53
	GPIO_GPPUD     = 0x0094 // GPIO Pin Pull-up/down Enable
	GPIO_GPPUDCLK0 = 0x0098 // GPIO Pin Pull-up/down Enable Clock 0
	GPIO_GPPUDCLK1 = 0x009c // GPIO Pin Pull-up/down Enable Clock 1
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	// Map logical pins to physical pins
	pinmap = map[uint]gopi.GPIOPin{
		3:  gopi.GPIOPin(2),  // I2C_SDA1      On Rpi1 A/B Rev1: GPIO0 otherwise GPIO2
		5:  gopi.GPIOPin(3),  // I2C_SCL1      On Rpi1 A/B Rev1: GPIO1 otherwise GPIO3
		7:  gopi.GPIOPin(4),  // GPIO_CLK
		8:  gopi.GPIOPin(14), // TXD0
		10: gopi.GPIOPin(15), // RXD0
		11: gopi.GPIOPin(17), // GPIO_GEN0
		12: gopi.GPIOPin(18), // GPIO_GEN1
		13: gopi.GPIOPin(27), // GPIO_GEN2     On Rpi1 A/B Rev1: GPIO21 otherwise GPIO27
		15: gopi.GPIOPin(22), // GPIO_GEN3
		16: gopi.GPIOPin(23), // GPIO_GEN4
		18: gopi.GPIOPin(24), // GPIO_GEN5
		19: gopi.GPIOPin(10), // SPI_MOSI
		21: gopi.GPIOPin(9),  // SPI_MOSO
		22: gopi.GPIOPin(25), // GPIO_GEN6
		23: gopi.GPIOPin(11), // SPI_CLK
		24: gopi.GPIOPin(8),  // SPI_CE0_N
		26: gopi.GPIOPin(7),  // SPI_CE1_N
		29: gopi.GPIOPin(5),  // Not on Rpi1 (all below here)
		31: gopi.GPIOPin(6),
		32: gopi.GPIOPin(12),
		33: gopi.GPIOPin(13),
		35: gopi.GPIOPin(19),
		36: gopi.GPIOPin(16),
		37: gopi.GPIOPin(26),
		38: gopi.GPIOPin(20),
		40: gopi.GPIOPin(21),
	}
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config GPIO) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.hw.rpi.GPIO.Open{ }")

	this := new(gpio)
	this.log = logger

	// Get product information
	if hw, ok := config.Hardware.(*hardware); ok == false || hw == nil {
		return nil, gopi.ErrBadParameter
	} else if product, err := hw.GetProduct(); err != nil {
		return nil, err
	} else {
		this.product = product
	}

	// Create pin mappings. Because there is some variation between the different
	// models and revisions, we use some logic in a private method to fudge
	this.pins = make(map[gopi.GPIOPin]uint)
	for pin := range pinmap {
		if logical := this.PhysicalPin(pin); logical != gopi.GPIO_PIN_NONE {
			this.pins[logical] = pin
		}
	}

	// Open the /dev/mem and provide offset & size for accessing memory
	if file, base, size, err := gpioOpenDevice(); err != nil {
		return nil, err
	} else {
		defer file.Close()

		// Lock memory
		this.memlock.Lock()
		defer this.memlock.Unlock()

		// Memory map GPIO registers to byte array
		if mem8, err := syscall.Mmap(int(file.Fd()), int64(base+GPIO_BASE), int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED); err != nil {
			return nil, err
		} else {
			this.mem8 = mem8
		}

		// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
		header := *(*reflect.SliceHeader)(unsafe.Pointer(&this.mem8))
		header.Len /= (32 / 8)
		header.Cap /= (32 / 8)
		this.mem32 = *(*[]uint32)(unsafe.Pointer(&header))

		// Success
		return this, nil
	}
}

// Close
func (this *gpio) Close() error {
	this.log.Debug("sys.hw.rpi.GPIO.Close{ }")

	// Unmap memory and return error
	this.memlock.Lock()
	defer this.memlock.Unlock()
	return syscall.Munmap(this.mem8)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *gpio) String() string {
	return fmt.Sprintf("sys.hw.rpi.GPIO{ physical_pins=%v logical_pins=%v product=%v }", this.NumberOfPhysicalPins(), this.pins, this.product)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENT INTERFACE

// NumberOfPhysicalPins returns number of physical pins
func (this *gpio) NumberOfPhysicalPins() uint {
	if this.product.model == RPI_MODEL_A || this.product.model == RPI_MODEL_B {
		return uint(26)
	} else {
		return uint(40)
	}
}

// Pins() returns array of available logical pins
func (this *gpio) Pins() []gopi.GPIOPin {
	pins := make([]gopi.GPIOPin, GPIO_MAXPINS)
	for i := 0; i < GPIO_MAXPINS; i++ {
		pins[i] = gopi.GPIOPin(i)
	}
	return pins
}

// PhysicalPin returns logical pin for physical pin number. Returns
// GPIO_PIN_NONE where there is no logical pin at that position
func (this *gpio) PhysicalPin(pin uint) gopi.GPIOPin {

	// Check for Raspberry Pi Version 1 and fudge things a little
	if this.product.model == RPI_MODEL_A || this.product.model == RPI_MODEL_B {
		// pin can be 1-28
		if pin < 1 || pin > 28 {
			return gopi.GPIO_PIN_NONE
		}
		if this.product.revision == Revision(1) && pin == 3 {
			return gopi.GPIOPin(0)
		}
		if this.product.revision == Revision(1) && pin == 5 {
			return gopi.GPIOPin(1)
		}
		if this.product.revision == Revision(1) && pin == 13 {
			return gopi.GPIOPin(21)
		}
	}

	// now do things normally...
	if logical_pin, ok := pinmap[pin]; ok == false {
		return gopi.GPIO_PIN_NONE
	} else {
		return logical_pin
	}
}

// PhysicalPinForPin returns physical pin number for logical pin.
// Returns 0 where there is no physical pin for this logical pin
func (this *gpio) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	if physical, ok := this.pins[logical]; ok == false {
		return 0
	} else {
		return physical
	}
}

// ReadPin reads pin state or returns LOW otherwise
func (this *gpio) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	var register uint32

	this.memlock.Lock()
	defer this.memlock.Unlock()

	if uint8(logical) <= uint8(31) {
		// GPIO0 - GPIO31
		register = this.mem32[GPIO_GPLVL0>>2]
	} else {
		// GPIO32 - GPIO53
		register = this.mem32[GPIO_GPLVL1>>2]
	}
	if (register & (1 << (uint8(logical) & 31))) != 0 {
		return gopi.GPIO_HIGH
	}
	return gopi.GPIO_LOW
}

// Write pin state
func (this *gpio) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {
	this.memlock.Lock()
	defer this.memlock.Unlock()

	v := uint32(1 << (uint8(logical) & 31))

	switch state {
	case gopi.GPIO_LOW:
		if uint8(logical) <= uint8(31) {
			this.mem32[GPIO_GPCLR0>>2] = v
		} else {
			this.mem32[GPIO_GPCLR1>>2] = v
		}
	case gopi.GPIO_HIGH:
		if uint8(logical) <= uint8(31) {
			this.mem32[GPIO_GPSET0>>2] = v
		} else {
			this.mem32[GPIO_GPSET1>>2] = v
		}
	}
}

// Get pin mode
func (this *gpio) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	// return the register and the number of bits to shift to
	// access the current mode
	register, shift := gopiPinToRegister(logical)

	this.memlock.Lock()
	defer this.memlock.Unlock()

	// Retrieve register, shift to the right, and return last three bits
	return gopi.GPIOMode((this.mem32[register>>2] >> shift) & 7)
}

// Set pin mode
func (this *gpio) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {
	// get register and the number of bits to shift to
	// access the current mode
	register, shift := gopiPinToRegister(logical)

	this.memlock.Lock()
	defer this.memlock.Unlock()

	this.mem32[register>>2] = (this.mem32[register>>2] &^ (7 << shift)) | (uint32(mode) << shift)
}

// Set pull mode
func (this *gpio) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) error {
	return gopi.ErrNotImplemented
}

// Watch is not implemented
func (this *gpio) Watch(gopi.GPIOPin, gopi.GPIOEdge) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PUBLISHER INTERFACE IS NOT IMPLEMENTED

func (this *gpio) Subscribe() <-chan gopi.Event {
	return nil
}

func (this *gpio) Unsubscribe(<-chan gopi.Event) {
	// Empty
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func gpioOpenDevice() (*os.File, uint32, uint32, error) {
	// open GPIO memory mapped file, or if that doesn't exist
	// attempt /dev/mem which would only work for root user
	if file, err := os.OpenFile(GPIO_DEV_GPIOMEM, os.O_RDWR|os.O_SYNC, 0); err == nil {
		return file, 0, GPIO_SIZE, nil
	} else if os.IsNotExist(err) {
		return nil, 0, 0, err
	} else if file, err = os.OpenFile(GPIO_DEV_MEM, os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, 0, 0, err
	} else {
		return file, bcmHostGetPeripheralAddress(), GPIO_SIZE, nil
	}
}

func gopiPinToRegister(pin gopi.GPIOPin) (uint, uint) {
	p := int(pin)
	switch {
	case p >= 0 && p <= 9:
		return GPIO_GPFSEL0, uint(p * 3)
	case p >= 10 && p <= 19:
		return GPIO_GPFSEL1, uint((p - 10) * 3)
	case p >= 20 && p <= 29:
		return GPIO_GPFSEL2, uint((p - 20) * 3)
	case p >= 30 && p <= 39:
		return GPIO_GPFSEL3, uint((p - 30) * 3)
	case p >= 40 && p <= 49:
		return GPIO_GPFSEL4, uint((p - 40) * 3)
	default:
		return GPIO_GPFSEL5, uint((p - 50) * 3)
	}
}
