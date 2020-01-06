# The Hardware Platform

In this section I describe how to interact with your hardware platform,
including:

 * Getting information about your hardware, including graphical display;
 * Interacting with the General Purpose Input-Output (GPIO) interface;
 * The I²C interface;
 * The SPI interface.

All of these features are available on the Raspberry Pi but not necessarily on other platforms.

## Hardware Information

The Platform Unit returns some information about the platform your tool is running on.

{% hint style="info" %}
| Parameter        | Value                   |
| ---------------- | ----------------------- |
| Name             | `gopi/platform`         |
| Interface        | `gopi.Platform`         |
| Type             | `gopi.UNIT_PLATFORM`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/platform` |
| Compatibility    | Linux, Darwin, Raspberry Pi |
{% endhint %}

Here is the interface which the platform unit adheres to:

```go
type gopi.Platform interface {
	gopi.Unit

	Product() string 	// Product returns product name
	Type() gopi.PlatformType 	// Type returns flags identifying platform type
	SerialNumber() string 	// SerialNumber returns unique serial number for host
	Uptime() time.Duration 	// Uptime returns uptime for host
	LoadAverages() (float64, float64, float64) 	// LoadAverages returns 1, 5 and 15 minute load averages
	NumberOfDisplays() uint 	// NumberOfDisplays returns the number of possible displays for this host
}
```

Here's an example of accessing the platform information in your `Main` function:

```go
func Main(app gopi.App, args []string) error {
	platform := app.Platform()
    fmt.Println(platform.Type(),platform.Product(),platform.SerialNumber())
    // ...
}
```

There are some platform differences with the information returned:

  * On Linux, the generic name `linux` is returned for product;
  * On Linux, a Mac Address is returned for the serial number;
  * On Darwin, the product is a product code (ie, "MacPro1,2") rather than name.
  * On Darwin and Linux, the number of displays is returned as zero as these platform displays are not yet supported.


## Displays

The Display Unit returns some information about your display. When importing
this unit into your tool, the command line flag `-display` can be used to choose
the display. An error will be returned when trying to use this unit on Linux or Darwin when the tool is run.

{% hint style="info" %}
| Parameter        | Value                   |
| ---------------- | ----------------------- |
| Name             | `gopi/display`         |
| Interface        | `gopi.Display`         |
| Type             | `gopi.UNIT_DISPLAY`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/display` |
| Compatibility    | Raspberry Pi |
{% endhint %}

Here's an example of returning information about the display:

```go
func Main(app gopi.App, args []string) error {
	display := app.Display()
    fmt.Println(display.Name(),display.Size())
    // ...
}
```

Here is the interface for a display:

```go
type gopi.Display interface {
    gopi.Unit

    DisplayId() uint 	// Return display number
	Name() string // Return name of the display
	Size() (uint32, uint32) // Return display size for nominated display number
	PixelsPerInch() uint32 // Return the PPI (pixels-per-inch) for the display
}
```

## I²C Interface

I²C is a serial protocol for two-wire interface to connect low-speed devices like sensors, A/D and D/A converters and other similar peripherals in embedded systems. It was invented by Philips and now it is used by almost all major IC manufacturers. For more information see [Wikipedia](https://en.wikipedia.org/wiki/I%C2%B2C).

The I²C unit allows you to read and write data with daisy-chained peripherals, each of which should
have a unique address.

{% hint style="info" %}
| Parameter        | Value               |
| ---------------- | ------------------- |
| Name             | `gopi/i2c`   |
| Interface        | `gopi.I2C`         |
| Type             | `gopi.UNIT_I2C`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/i2c` |
| Compatibility    | Linux               |
{% endhint %}

The unit adheres to the following interface:

```go
type gopi.I2C interface {
	gopi.Unit

	SetSlave(uint8) error
	GetSlave() uint8
	DetectSlave(uint8) (bool, error)

	// Read 
	ReadUint8(reg uint8) (uint8, error)
	ReadInt8(reg uint8) (int8, error)
	ReadUint16(reg uint8) (uint16, error)
	ReadInt16(reg uint8) (int16, error)
	ReadBlock(reg, length uint8) ([]byte, error)

	// Write
	WriteUint8(reg, value uint8) error
	WriteInt8(reg uint8, value int8) error
	WriteUint16(reg uint8, value uint16) error
	WriteInt16(reg uint8, value int16) error
}
```

You need to set a slave address when using the tool, which is a value between `0x00` and `0x7F`. You can use the `DetectSlave` method which
returns `true` if a peripheral was found at a particular slave address. For example,

```go
func Main(app gopi.App, args []string) error {
	i2c := app.I2C()
	slave := app.Flags().GetUint("slave",gopi.FLAG_NS_DEFAULT)
	if detected, err := i2c.DetectSlave(slave); detected == false {
		return fmt.Errorf("No peripheral detected")
	} else if err := this.i2c.SetSlave(slave); err != nil {
		return err
	} else if reg0, err := this.i2c.ReadInt16(0) {
		fmt.Println("REG0=",reg0)		
	}
    // ...
}
```

The unit adds an additional commmand line flag of `-i2c.bus` to
select which interface to attach to. On the Raspberry Pi, you need to enable the interface using the `raspi-config` command and ensure
your user has the correct permissions to access the device using the 
following command:

```bash
bash% sudo usermod -a -G i2c ${USER}
```

There's more information about enabling it [here](https://www.electronicwings.com/raspberry-pi/raspberry-pi-i2c).

There are some examples of using the I2C unit in the [sensors](github.com/djthorpe/sensors) repository
including temperature, light and humidity measurement using
I²C peripherals.

## SPI Interface

The Serial Peripheral Interface (SPI) is a synchronous serial communication interface for embedded systems. More information is
available on [Wikipedia](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface).

The SPI Unit allows you to read and write data, or do bi-directional
transfers. In order to use the Unit, here are the parameters:

{% hint style="info" %}
| Parameter        | Value               |
| ---------------- | ------------------- |
| Name             | `gopi/spi`   |
| Interface        | `gopi.SPI`         |
| Type             | `gopi.UNIT_SPI`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/spi` |
| Compatibility    | Linux               |
{% endhint %}

The unit adheres to the following interface:

```go

// SPI implements the SPI interface for sensors, etc.
type gopi.SPI interface {
	gopi.Unit

	Mode() gopi.SPIMode
	MaxSpeedHz() uint32
	BitsPerWord() uint8

	SetMode(gopi.SPIMode) error
	SetMaxSpeedHz(uint32) error
	SetBitsPerWord(uint8) error

	Read(len uint32) ([]byte, error)
	Write(send []byte) error
	Transfer(send []byte) ([]byte, error)
}
```

The unit adds the flags `-spi.bus` and `-spi.slave` to the
command-line flags in order to select the correct device.

On the Raspberry Pi, you need to enable the interface using the `raspi-config` command and ensure
your user has the correct permissions to access the device using the 
following command:

```bash
bash% sudo usermod -a -G spi ${USER}
```

There's more information about the Raspberry Pi implementation [here](https://www.raspberrypi.org/documentation/hardware/raspberrypi/spi/README.md).

## GPIO Interface

_Documentation to be written_

{% hint style="info" %}
| Parameter        | Value               |
| ---------------- | ------------------- |
| Name             | `gopi/gpio/linux`   |
| Or               | `gopi/gpio/rpi`     |
| Interface        | `gopi.GPIO`         |
| Type             | `gopi.UNIT_GPIO`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/gpio` |
| Events           | `gopi.GPIOEvent`    |
| Compatibility    | Linux, Raspberry Pi    |
{% endhint %}

