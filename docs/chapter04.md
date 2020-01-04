# The Hardware Platform

In this section I describe how to interact with your hardware platform,
including:

 * Getting information about your hardware, including graphical display;
 * Interacting with the General Purpose Input-Output (GPIO) interface;
 * The I2C interface;
 * The SPI interface.

All of these features are available on the Raspberry Pi but not necessarily on other platforms.

## The Hardware Platform Unit

| Parameter        | Value                   |
| ---------------- | ----------------------- |
| Name             | `gopi/platform`         |
| Interface        | `gopi.Platform`         |
| Type             | `gopi.UNIT_PLATFORM`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/platform` |
| Compatibility    | Linux, Darwin, Raspberry Pi |


## The Display Unit

| Parameter        | Value                   |
| ---------------- | ----------------------- |
| Name             | `gopi/display`         |
| Interface        | `gopi.Display`         |
| Type             | `gopi.UNIT_DISPLAY`    |
| Import           | `github.com/djthorpe/gopi/v2/unit/display` |
| Compatibility    | Raspberry Pi |

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

