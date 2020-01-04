# The Hardware Platform

In this section I describe how to interact with your hardware platform,
including:

 * Getting information about your hardware, including graphical display;
 * Interacting with the General Purpose Input-Output (GPIO) interface;
 * The I2C interface;
 * The SPI interface.

All of these features are available on the Raspberry Pi but not necessarily on other platforms.

## The Hardware Platform Unit

The Platform Unit returns some information about the platform your tool is running on.

| Parameter        | Value                   |
| ---------------- | ----------------------- |
| Name             | `gopi/platform`         |
| Interface        | `gopi.Platform`         |
| Type             | `gopi.UNIT_PLATFORM`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/platform` |
| Compatibility    | Linux, Darwin, Raspberry Pi |

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


## The Display Unit

The Display Unit returns some information about your display. When importing
this unit into your tool, the command line flag `-display` can be used to choose
the display. An error will be returned when trying to use this unit on Linux or Darwin when the tool is run.

| Parameter        | Value                   |
| ---------------- | ----------------------- |
| Name             | `gopi/display`         |
| Interface        | `gopi.Display`         |
| Type             | `gopi.UNIT_DISPLAY`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/display` |
| Compatibility    | Raspberry Pi |


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
type Display interface {
    gopi.Unit

    DisplayId() uint 	// Return display number
	Name() string // Return name of the display
	Size() (uint32, uint32) // Return display size for nominated display number
	PixelsPerInch() uint32 // Return the PPI (pixels-per-inch) for the display
}
```

## The GPIO Unit

| Parameter        | Value               |
| ---------------- | ------------------- |
| Name             | `gopi/gpio/linux`   |
| Or               | `gopi/gpio/rpi`     |
| Interface        | `gopi.GPIO`         |
| Type             | `gopi.UNIT_GPIO`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/gpio` |
| Events           | `gopi.GPIOEvent`    |
| Compatibility    | Linux, Raspberry Pi    |

## The I2C Unit

| Parameter        | Value               |
| ---------------- | ------------------- |
| Name             | `gopi/i2c`   |
| Interface        | `gopi.I2C`         |
| Type             | `gopi.UNIT_I2C`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/i2c` |
| Compatibility    | Linux               |


## The SPI Unit

| Parameter        | Value               |
| ---------------- | ------------------- |
| Name             | `gopi/spi`   |
| Interface        | `gopi.SPI`         |
| Type             | `gopi.UNIT_SPI`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/spi` |
| Compatibility    | Linux               |

