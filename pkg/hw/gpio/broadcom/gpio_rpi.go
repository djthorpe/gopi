// +build rpi

package broadcom

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/hw/gpio"
	"github.com/djthorpe/gopi/v3/pkg/hw/platform"
	"github.com/djthorpe/gopi/v3/pkg/sys/rpi"
)

type GPIO struct {
	gopi.Unit
	sync.RWMutex
	gopi.Publisher
	*platform.Platform

	product rpi.ProductInfo
	pins    map[gopi.GPIOPin]uint           // map of logical to physical pins
	mem8    []uint8                         // access GPIO as bytes
	mem32   []uint32                        // access GPIO as uint32
	watch   map[gopi.GPIOPin]gopi.GPIOState // current pin state
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
	watchDelta = 250 * time.Millisecond // Updates pin state every 250ms
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
// IMPLEMENTATION

func (this *GPIO) New(gopi.Config) error {

	if _, product, err := rpi.VCGetSerialProduct(); err != nil {
		return err
	} else {
		this.product = rpi.NewProductInfo(product)
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
		return err
	} else {
		defer file.Close()

		// Memory map GPIO registers to byte array
		if mem8, err := syscall.Mmap(int(file.Fd()), int64(base+GPIO_BASE), int(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED); err != nil {
			return err
		} else {
			this.mem8 = mem8
		}

		// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
		header := *(*reflect.SliceHeader)(unsafe.Pointer(&this.mem8))
		header.Len /= (32 / 8)
		header.Cap /= (32 / 8)
		this.mem32 = *(*[]uint32)(unsafe.Pointer(&header))
	}

	// Check length of arrays
	if len(this.mem8) == 0 || len(this.mem32) == 0 {
		return gopi.ErrInternalAppError.WithPrefix("New")
	}

	// Set up pin watching
	this.watch = make(map[gopi.GPIOPin]gopi.GPIOState)

	// Return success
	return nil
}

func (this *GPIO) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if err := syscall.Munmap(this.mem8); err != nil {
		return os.NewSyscallError("munmap", err)
	}

	// Release resources
	this.pins = nil
	this.watch = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *GPIO) Run(ctx context.Context) error {
	timer := time.NewTicker(watchDelta)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			this.changeWatchState()
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *GPIO) String() string {
	str := "<gpio.broadcom"
	if p := this.NumberOfPhysicalPins(); p > 0 {
		str += " number_of_physical_pins=" + fmt.Sprint(p)
	}
	if l := this.pins; len(l) > 0 {
		str += " pins=" + fmt.Sprint(l)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PINS

// Return number of physical pins, or 0 if if cannot be returned
// or nothing is known about physical pins
func (this *GPIO) NumberOfPhysicalPins() uint {
	if this.product.Model == rpi.RPI_MODEL_A || this.product.Model == rpi.RPI_MODEL_B {
		return uint(26)
	} else {
		return uint(40)
	}
}

// Return array of available logical pins or nil if nothing is
// known about pins
func (this *GPIO) Pins() []gopi.GPIOPin {
	pins := make([]gopi.GPIOPin, GPIO_MAXPINS)
	for i := 0; i < GPIO_MAXPINS; i++ {
		pins[i] = gopi.GPIOPin(i)
	}
	return pins
}

// Return logical pin for physical pin number. Returns
// GPIO_PIN_NONE where there is no logical pin at that position
// or we don't know about the physical pins
func (this *GPIO) PhysicalPin(pin uint) gopi.GPIOPin {
	// Check for Raspberry Pi Version 1 and fudge things a little
	if this.product.Model == rpi.RPI_MODEL_A || this.product.Model == rpi.RPI_MODEL_B {
		// pin can be 1-28
		if pin < 1 || pin > 28 {
			return gopi.GPIO_PIN_NONE
		}
		if this.product.Revision == rpi.Revision(1) && pin == 3 {
			return gopi.GPIOPin(0)
		}
		if this.product.Revision == rpi.Revision(1) && pin == 5 {
			return gopi.GPIOPin(1)
		}
		if this.product.Revision == rpi.Revision(1) && pin == 13 {
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

// Return physical pin number for logical pin. Returns 0 where there
// is no physical pin for this logical pin, or we don't know anything
// about the layout
func (this *GPIO) PhysicalPinForPin(logical gopi.GPIOPin) uint {
	if physical, ok := this.pins[logical]; ok == false {
		return 0
	} else {
		return physical
	}
}

// Read pin state
func (this *GPIO) ReadPin(logical gopi.GPIOPin) gopi.GPIOState {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var register uint32
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
func (this *GPIO) WritePin(logical gopi.GPIOPin, state gopi.GPIOState) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	value := uint32(1 << (uint8(logical) & 31))
	switch state {
	case gopi.GPIO_LOW:
		if uint8(logical) <= uint8(31) {
			this.mem32[GPIO_GPCLR0>>2] = value
		} else {
			this.mem32[GPIO_GPCLR1>>2] = value
		}
	case gopi.GPIO_HIGH:
		if uint8(logical) <= uint8(31) {
			this.mem32[GPIO_GPSET0>>2] = value
		} else {
			this.mem32[GPIO_GPSET1>>2] = value
		}
	}
}

// Get pin mode
func (this *GPIO) GetPinMode(logical gopi.GPIOPin) gopi.GPIOMode {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// return the register and the number of bits to shift to
	// access the current mode
	register, shift := gopiPinToRegister(logical)

	// Retrieve register, shift to the right, and return last three bits
	return gopi.GPIOMode((this.mem32[register>>2] >> shift) & 7)
}

// Set pin mode
func (this *GPIO) SetPinMode(logical gopi.GPIOPin, mode gopi.GPIOMode) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// get register and the number of bits to shift to
	// access the current mode
	register, shift := gopiPinToRegister(logical)

	// Set register
	this.mem32[register>>2] = (this.mem32[register>>2] &^ (7 << shift)) | (uint32(mode) << shift)
}

// Set pull mode to pull down or pull up - will
// return ErrNotImplemented if not supported
func (this *GPIO) SetPullMode(logical gopi.GPIOPin, pull gopi.GPIOPull) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check pin to make sure there is a physical pin mapping
	if this.PhysicalPinForPin(logical) == 0 {
		return gopi.ErrBadParameter.WithPrefix(fmt.Sprint(logical))
	}

	// Set the low two bits of register to 0 (off) 1 (down) or 2 (up)
	switch pull {
	case gopi.GPIO_PULL_UP, gopi.GPIO_PULL_DOWN:
		this.mem32[GPIO_GPPUD] |= uint32(pull)
	case gopi.GPIO_PULL_OFF:
		this.mem32[GPIO_GPPUD] &^= 3
	}

	// Wait for 150 cycles
	time.Sleep(time.Microsecond)

	// Determine clock register
	clockReg := GPIO_GPPUDCLK0
	if logical >= gopi.GPIOPin(32) {
		clockReg = GPIO_GPPUDCLK1
	}

	// Clock it in
	this.mem32[clockReg] = 1 << (logical % 32)

	// Wait for value to clock in
	time.Sleep(time.Microsecond)

	// Write 00 to the register to clear it
	this.mem32[GPIO_GPPUD] &^= 3

	// Wait for value to clock in
	time.Sleep(time.Microsecond)

	// Remove the clock
	this.mem32[clockReg] = 0

	// Return success
	return nil
}

// Start watching for rising and/or falling edge,
// or stop watching when GPIO_EDGE_NONE is passed.
// Will return ErrNotImplemented if not supported
func (this *GPIO) Watch(pin gopi.GPIOPin, edge gopi.GPIOEdge) error {
	// Check pin mode is INPUT
	if mode := this.GetPinMode(pin); mode != gopi.GPIO_INPUT {
		return gopi.ErrOutOfOrder.WithPrefix("Watch", pin)
	}

	// Get existing state of pin
	state := gopi.GPIO_LOW
	if edge != gopi.GPIO_EDGE_NONE {
		state = this.ReadPin(pin)
	}

	// Lock for writing
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Delete watch or set existing state
	if edge == gopi.GPIO_EDGE_NONE {
		delete(this.watch, pin)
		return nil
	} else {
		this.watch[pin] = state
	}

	// Return success
	return nil
}

func (this *GPIO) changeWatchState() {
	for pin, state := range this.watch {
		if newstate := this.ReadPin(pin); newstate == state {
			continue
		} else {
			this.RWMutex.Lock()
			defer this.RWMutex.Unlock()
			this.watch[pin] = newstate
		}
		if this.Publisher != nil {
			edge := gopi.GPIO_EDGE_NONE
			if state == gopi.GPIO_LOW {
				edge = gopi.GPIO_EDGE_RISING
			} else {
				edge = gopi.GPIO_EDGE_FALLING
			}
			this.Publisher.Emit(gpio.NewEvent(fmt.Sprint(pin), pin, edge), true)
		}
	}
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
		return file, rpi.BCMHostGetPeripheralAddress(), GPIO_SIZE, nil
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
