
## Interfacing with GPIO, I²C and SPI

The Raspberry Pi includes a general-purpose input/output (GPIO)
header which can be used to interface with a variety of hardware
like sensors, switches, lights, motors, and so forth. The gopi
framework includes modules which allow you to easily control
the devices through GPIO, I²C and SPI protocols.

There are two GPIO variants and implementations for I²C and SPI
using Linux kernel drivers. There is also an LIRC module which
implements the infra-red sending and receiving for remote
controls:

| Name        | Use                 | Abstract Interface    | Import                                      |
| -- | -- | -- | -- |
| "rpi/gpio"       | app.GPIO          | `gopi.GPIO`         | `github.com/djthorpe/gopi/sys/hw/rpi`      |
| "linux/gpio"     | app.GPIO          | `gopi.GPIO`         | `github.com/djthorpe/gopi/sys/hw/linux`    |
| "linux/spi"      | app.SPI           | `gopi.SPI`          | `github.com/djthorpe/gopi/sys/hw/linux`    |
| "linux/i2c"      | app.I2C           | `gopi.I2C`          | `github.com/djthorpe/gopi/sys/hw/linux`    |
| "linux/lirc"     | app.LIRC          | `gopi.LIRC`         | `github.com/djthorpe/gopi/sys/hw/linux`    |


### The GPIO interface

The GPIO module can be included either in a Linux variant or
a Raspberry Pi variant. The Linux variant can emit events when
a GPIO input pin is changed from low to high state (rising edge) or
high to low state (falling edge). The Raspberry Pi module can't do that,
but it can set the mode on the pin and is generally faster since
it uses memory-mapped registers.

Here is the interface for GPIO:

```
type GPIO interface {
	Driver
	Publisher

	// Return number of physical pins, or 0 if if cannot be returned
	// or nothing is known about physical pins
	NumberOfPhysicalPins() uint

	// Return array of available logical pins or nil if nothing is
	// known about pins
	Pins() []GPIOPin

	// Return logical pin for physical pin number. Returns
	// GPIO_PIN_NONE where there is no logical pin at that position
	// or we don't now about the physical pins
	PhysicalPin(uint) GPIOPin

	// Return physical pin number for logical pin. Returns 0 where there
	// is no physical pin for this logical pin, or we don't know anything
	// about the layout
	PhysicalPinForPin(GPIOPin) uint

	// Read pin state
	ReadPin(GPIOPin) GPIOState

	// Write pin state
	WritePin(GPIOPin, GPIOState)

	// Get pin mode
	GetPinMode(GPIOPin) GPIOMode

	// Set pin mode
	SetPinMode(GPIOPin, GPIOMode)

	// Set pull mode to pull down or pull up - will
	// return ErrNotImplemented if not supported
	SetPullMode(GPIOPin, GPIOPull) error

	// Start watching for rising and/or falling edge,
	// or stop watching when GPIO_EDGE_NONE is passed.
	// Will return ErrNotImplemented if not supported
	Watch(GPIOPin, GPIOEdge) error
}
```

The first concept to take note of is the difference between the
physical pins and the logical pin numbers. On the Raspberry Pi,
the mapping between these can vary between models. The
physical pins are referred to by unsigned integers and the
logical pins as type `GPIOPin`. You can map between them using
the methods `PhysicalPin` to convert a physical pin number into
a logical GPIOPin and `PhysicalPinForPin` to do the reverse.

In order to make a pin high (or vice versa) you could implement
your main function as follows:

```
func Main(app *gopi.AppInstance,done chan<- struct{}) error {
	pin, _ := app.AppFlags.GetUint("pin")
	gpio_pin := app.GPIO.PhysicalPin(pin)
	app.GPIO.SetPinMode(gpio_pin,gopi.GPIO_OUTPUT)
	app.GPIO.WritePin(gpio_pin,gopi.GPIO_HIGH)
	return nil
}
```

On the Raspberry Pi, you can configure your header so that each
physical pin can be an input pin, output pin or one of several
functions [defined here](). The pin modes are as follows:

| Name        | Mode       |
| -- | -- |
| GPIO_INPUT  | Input Pin            |
| GPIO_OUTPUT | Output Pin           |
| GPIO_ALT0   | Alternate Function 0 |
| GPIO_ALT1   | Alternate Function 1 |
| GPIO_ALT2   | Alternate Function 2 |
| GPIO_ALT3   | Alternate Function 3 |
| GPIO_ALT4   | Alternate Function 4 |
| GPIO_ALT5   | Alternate Function 5 |

When using the Linux driver, it's possible to watch input pins so that
when the pin transitions from low to high or vice-versa, an even is emitted:

```
type GPIOEvent interface {
	Event

	// Pin returns the pin on which the event occurred
	Pin() GPIOPin

	// Edge returns whether the pin value is rising or falling
	// or will return NONE if not defined
	Edge() GPIOEdge
}
```

In order to receive events, you would need to `Subscribe` to them and
respond to emitted events in the background. For example,

```
TODO
```


### The I²C interface

type I2C interface {
	Driver

	// Set current slave address
	SetSlave(uint8) error

	// Get current slave address
	GetSlave() uint8

	// Return true if a slave was detected at a particular address
	DetectSlave(uint8) (bool, error)

	// Read Byte (8-bits), Word (16-bits) & Block ([]byte) from registers
	ReadUint8(reg uint8) (uint8, error)
	ReadInt8(reg uint8) (int8, error)
	ReadUint16(reg uint8) (uint16, error)
	ReadInt16(reg uint8) (int16, error)
	ReadBlock(reg, length uint8) ([]byte, error)

	// Write Byte (8-bits) & Word (16-bits) to registers
	WriteUint8(reg, value uint8) error
	WriteInt8(reg uint8, value int8) error
	WriteUint16(reg uint8, value uint16) error
	WriteInt16(reg uint8, value int16) error
}


### The SPI interface

type SPI interface {
	Driver

	// Get SPI mode
	Mode() SPIMode
	// Get SPI speed
	MaxSpeedHz() uint32
	// Get Bits Per Word
	BitsPerWord() uint8
	// Set SPI mode
	SetMode(SPIMode) error
	// Set SPI speed
	SetMaxSpeedHz(uint32) error
	// Set Bits Per Word
	SetBitsPerWord(uint8) error

	// Read/Write
	Transfer(send []byte) ([]byte, error)

	// Read
	Read(len uint32) ([]byte, error)

	// Write
	Write(send []byte) error
}


### The LIRC interface

type LIRC interface {
	Driver
	Publisher

	// Get receive and send modes
	RcvMode() LIRCMode
	SendMode() LIRCMode
	SetRcvMode(mode LIRCMode) error
	SetSendMode(mode LIRCMode) error

	// Receive parameters
	GetRcvResolution() (uint32, error)
	SetRcvTimeout(micros uint32) error
	SetRcvTimeoutReports(enable bool) error
	SetRcvCarrierHz(value uint32) error
	SetRcvCarrierRangeHz(min uint32, max uint32) error

	// Send parameters
	SetSendCarrierHz(value uint32) error
	SetSendDutyCycle(value uint32) error

	// Send Pulse Mode, values are in milliseconds
	PulseSend(values []uint32) error
}
