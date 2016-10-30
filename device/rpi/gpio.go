/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"reflect"
	"unsafe"
)

import (
	gopi "../.."      /* import "github.com/djthorpe/gopi" */
	util "../../util" /* import "github.com/djthorpe/gopi/util" */
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type GPIO struct {
	Device  gopi.HardwareDriver
}

type GPIODriver struct {
	log      *util.LoggerDevice // logger
	memlock sync.Mutex
	model   Model // Device model
	revision PCBRevision // PCB revision
	mem8    []uint8  // access GPIO as bytes
	mem32   []uint32 // access GPIO as uint32
	pins    map[gopi.GPIOPin]uint // map of logical to physical pins
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
	GPIO_GPLVL0 = 0x0034 // Register to read pins GPIO0-GPIO31
	GPIO_GPLVL1 = 0x0038 // Register to read pins GPIO32-GPIO53
	GPIO_GPSET0 = 0x001C // Register to write HIGH to pins GPIO0-GPIO31
	GPIO_GPSET1 = 0x0020 // Register to write HIGH to pins GPIO32-GPIO53
	GPIO_GPCLR0 = 0x0028 // Register to write LOW to pins GPIO0-GPIO31
	GPIO_GPCLR1 = 0x002C // Register to write LOW to pins GPIO32-GPIO53
	GPIO_GPFSEL0 = 0x0000 // Pin modes for GPIO0-GPIO9
	GPIO_GPFSEL1 = 0x0004 // Pin modes for GPIO10-GPIO19
	GPIO_GPFSEL2 = 0x0008 // Pin modes for GPIO20-GPIO29
	GPIO_GPFSEL3 = 0x000C // Pin modes for GPIO30-GPIO39
	GPIO_GPFSEL4 = 0x0010 // Pin modes for GPIO40-GPIO49
	GPIO_GPFSEL5 = 0x0014 // Pin modes for GPIO50-GPIO53
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	// Map logical pins to physical pins
	pinmap = map[uint]gopi.GPIOPin{
		3: gopi.GPIOPin(2),     // I2C_SDA1      On Rpi1 A/B Rev1: GPIO0 otherwise GPIO2
		5: gopi.GPIOPin(3),     // I2C_SCL1      On Rpi1 A/B Rev1: GPIO1 otherwise GPIO3
		7: gopi.GPIOPin(4),     // GPIO_CLK
		8: gopi.GPIOPin(14),    // TXD0
		10: gopi.GPIOPin(15),   // RXD0
		11: gopi.GPIOPin(17),   // GPIO_GEN0
		12: gopi.GPIOPin(18),   // GPIO_GEN1
		13: gopi.GPIOPin(27),   // GPIO_GEN2     On Rpi1 A/B Rev1: GPIO21 otherwise GPIO27
		15: gopi.GPIOPin(22),   // GPIO_GEN3
		16: gopi.GPIOPin(23),   // GPIO_GEN4
		18: gopi.GPIOPin(24),   // GPIO_GEN5
		19: gopi.GPIOPin(10),   // SPI_MOSI
		21: gopi.GPIOPin(9),    // SPI_MOSO
		22: gopi.GPIOPin(25),   // GPIO_GEN6
		23: gopi.GPIOPin(11),   // SPI_CLK
		24: gopi.GPIOPin(8),    // SPI_CE0_N
		26: gopi.GPIOPin(7),    // SPI_CE1_N
		29: gopi.GPIOPin(5),    // Not on Rpi1 (all below here)
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

// Create new Display object, returns error if not possible
func (config GPIO) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	var err error
	log.Debug("<rpi.GPIO>Open")

	// create new GPIO driver
	this := new(GPIODriver)

	// Set logging & device
	this.log = log

	// Get Respberry Pi Model and Revision
	this.model, this.revision, err = config.Device.(*DeviceState).GetModel()
	if err != nil {
		return nil, err
	}

	// Create pin mappings. Because there is some variation between the different
	// models and revisions, we use some logic in a private method to fudge
	this.pins = make(map[gopi.GPIOPin]uint)
	for pin,_ := range pinmap {
		if logical := this.PhysicalPin(pin); logical != gopi.GPIO_PIN_NONE {
			this.pins[logical] = pin
		}
	}

	// Open the /dev/mem and provide offset & size for accessing memory
	file, peripheral_base, peripheral_size, err := gpioOpenDevice(config.Device.(*DeviceState))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Lock memory
	this.memlock.Lock()
	defer this.memlock.Unlock()

	// Memory map GPIO registers to byte array
	this.mem8, err = syscall.Mmap(int(file.Fd()), int64(peripheral_base + GPIO_BASE), int(peripheral_size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)
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

// Close GPIO connection
func (this *GPIODriver) Close() error {
	this.log.Debug("<rpi.GPIO>Close")

	// Unmap memory and return error
	this.memlock.Lock()
	defer this.memlock.Unlock()
	return syscall.Munmap(this.mem8)
}

// Strinfigy GPIO
func (this *GPIODriver) String() string {
	return fmt.Sprintf("<rpi.GOPI>{ %v }",this.pins)
}

// Return logical pin for physical pin number. Will return
// gopi.GPIO_PIN_NONE if there is no logical pin on that physical
// one
func (this *GPIODriver) PhysicalPin(pin uint) gopi.GPIOPin {

	// Check for Raspberry Pi Version 1 and fudge things a little
	if this.model == RPI_MODEL_A || this.model == RPI_MODEL_B {
		// pin can be 1-28
		if pin < 1 || pin > 28 {
			return gopi.GPIO_PIN_NONE
		}
		if this.revision == PCBRevision(1) && pin == 3 {
			return gopi.GPIOPin(0)
		}
		if this.revision == PCBRevision(1) && pin == 5 {
			return gopi.GPIOPin(1)
		}
		if this.revision == PCBRevision(1) && pin == 13 {
			return gopi.GPIOPin(21)
		}
	}

	// now do things normally...
	logical_pin, ok := pinmap[pin]
	if ok == false {
		return gopi.GPIO_PIN_NONE
	}
	return logical_pin
}

// Return physical pin number for logical pin, or 0 if there is no
// mapping from the logical pin to the physical one
func (this *GPIODriver) PhysicalPinForPin(pin gopi.GPIOPin) uint {
	physical, ok := this.pins[pin]
	if ok != true {
		return 0
	} else {
		return physical
	}
}

// Return all logical pins - some won't exist on the physical board
// so use the PhysicalPinForPin function in order to determine if
// these pins exist
func (this *GPIODriver) Pins() []gopi.GPIOPin {
	pins := make([]gopi.GPIOPin,GPIO_MAXPINS)
	for i := 0; i < GPIO_MAXPINS; i++ {
		pins[i] = gopi.GPIOPin(i)
	}
	return pins
}

////////////////////////////////////////////////////////////////////////////////

// ReadPin reads the state of a pin
func (this *GPIODriver) ReadPin(pin gopi.GPIOPin) gopi.GPIOState {
	var register uint32

	this.memlock.Lock()
	defer this.memlock.Unlock()

	if uint8(pin) <= uint8(31) {
		// GPIO0 - GPIO31
		register = this.mem32[GPIO_GPLVL0 >> 2]
	} else {
		// GPIO32 - GPIO53
		register = this.mem32[GPIO_GPLVL1 >> 2]
	}
	if (register & (1 << (uint8(pin) & 31))) != 0 {
		return gopi.GPIO_HIGH
	}
	return gopi.GPIO_LOW
}

// WritePin writes the state of a pin
func (this *GPIODriver) WritePin(pin gopi.GPIOPin,state gopi.GPIOState) {

	this.memlock.Lock()
	defer this.memlock.Unlock()

	v := uint32(1 << (uint8(pin) & 31))

	switch(state) {
	case gopi.GPIO_LOW:
		if uint8(pin) <= uint8(31) {
			this.mem32[GPIO_GPCLR0 >> 2] = v
		} else {
			this.mem32[GPIO_GPCLR1 >> 2] = v
		}
	case gopi.GPIO_HIGH:
		if uint8(pin) <= uint8(31) {
			this.mem32[GPIO_GPSET0 >> 2] = v
		} else {
			this.mem32[GPIO_GPSET1 >> 2] = v
		}
	}
}

// GetPinMode reads the current pin mode
func (this *GPIODriver) GetPinMode(pin gopi.GPIOPin) gopi.GPIOMode {
	// return the register and the number of bits to shift to
	// access the current mode
	register,shift := gopiPinToRegister(pin)

	this.memlock.Lock()
	defer this.memlock.Unlock()

	// Retrieve register, shift to the right, and return last three bits
	return gopi.GPIOMode((this.mem32[register >> 2] >> shift) & 7)
}

// SetPinMode writes the current pin mode
func (this *GPIODriver) SetPinMode(pin gopi.GPIOPin,mode gopi.GPIOMode) {
	// get register and the number of bits to shift to
	// access the current mode
	register,shift := gopiPinToRegister(pin)

	this.memlock.Lock()
	defer this.memlock.Unlock()

	this.mem32[register >> 2] = (this.mem32[register >> 2] &^ (7 << shift)) | (uint32(mode) << shift)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func gpioOpenDevice(device *DeviceState) (*os.File, uint32, uint32, error) {
	var file *os.File
	var err error
	var peripheral_base uint32

	// open GPIO file
	if file, err = os.OpenFile(GPIO_DEV_GPIOMEM, os.O_RDWR|os.O_SYNC, 0); os.IsNotExist(err) {
		file, err = os.OpenFile(GPIO_DEV_MEM, os.O_RDWR|os.O_SYNC, 0)
		if err != nil {
			return nil, 0, 0, err
		}
		// peripheral_base is not zero-based
		peripheral_base = device.GetPeripheralAddress()
	}
	if err != nil {
		return nil, 0, 0, err
	}
	return file, peripheral_base, GPIO_SIZE, nil
}

func gopiPinToRegister(pin gopi.GPIOPin) (uint,uint) {
	p := int(pin)
	switch {
	case p >= 0 && p <= 9:
		return GPIO_GPFSEL0,uint(p * 3)
	case p >= 10 && p <= 19:
		return GPIO_GPFSEL1,uint((p-10) * 3)
	case p >= 20 && p <= 29:
		return GPIO_GPFSEL2,uint((p-20) * 3)
	case p >= 30 && p <= 39:
		return GPIO_GPFSEL3,uint((p-30) * 3)
	case p >= 40 && p <= 49:
		return GPIO_GPFSEL4,uint((p-40) * 3)
	default:
		return GPIO_GPFSEL5,uint((p-50) * 3)
	}
}



