
# Communicating with devices on the I2C and SPI busses

The Raspberry Pi contains an I2C ("I squared C") and SPI ("Serial Peripheral Interface")
busses for communicating with devices which use these standards.

The following sections explain how to communicate with devices over these busses.

## Abstract interfaces

There are separate drivers for the I2C and SPI busses, but you communicate to
devices in similar ways. With the I2C bus, only two signal lines are required,
the data line (SDA) for bi-directional communication and the clock line (SCL)
for managing timing. You may have up to 128 slave devices on each bus and 
there are two busses on the Raspberry Pi.

With the SPI bus, four signal lines are required, the slave input line (MOSI),
the slave output line (MISO), the clock line (CLK) and the slave selection (CE)
line. Out of the box, only two slave devices can be attached to one SPI bus.
The Raspberry Pi has a single bus.

Here are the abstract interfaces for the drivers:

| **Import** | `github.com/djthorpe/gopi/hw` |
| -- | -- | -- |
| **Interface** | `hw.I2CDriver` | gopi.Driver, the I2C driver |
| **Interface** | `hw.SPIDriver` | gopi.Driver, the SPI driver |

## Concrete implementations

The concrete implementations for the drivers are implemented by communicating
with Linux via device drivers. Here are the configuration structs which are
used to create the drivers:

| **Import** | `github.com/djthorpe/device/linux` |
| -- | -- | -- |
| **Struct** | `linux.I2C` | Configuration for the I2C concrete Linux driver |
| **Struct** | `linux.SPI` | Configuration for the SPI concrete Linux driver |


