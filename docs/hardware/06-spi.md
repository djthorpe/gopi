# SPI interface

The Serial Peripheral Interface \(SPI\) is a synchronous serial communication interface for embedded systems. More information is available on [Wikipedia](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface).

The SPI Unit allows you to read and write data, or do bi-directional transfers. In order to use the Unit, here are the parameters:

{% hint style="info" %}
| Parameter | Value |
| :--- | :--- |
| Name | `gopi/spi` |
| Interface | `gopi.SPI` |
| Type | `gopi.UNIT_SPI` |
| Import | `github.com/djthorpe/gopi/v2/unit/spi` |
| Compatibility | Linux |
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

The unit adds the flags `-spi.bus` and `-spi.slave` to the command-line flags in order to select the correct device.

On the Raspberry Pi, you need to enable the interface using the `raspi-config` command and ensure your user has the correct permissions to access the device using the following command:

```bash
bash% sudo usermod -a -G spi ${USER}
```

There's more information about the Raspberry Pi implementation [here](https://www.raspberrypi.org/documentation/hardware/raspberrypi/spi/README.md).

