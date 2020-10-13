package waveshare

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	gpiobcm "github.com/djthorpe/gopi/v3/pkg/hw/gpiobcm"
	spi "github.com/djthorpe/gopi/v3/pkg/hw/spi"
)

type EPD struct {
	gopi.Unit
	*gpiobcm.GPIO
	*spi.Devices

	bus, slave *uint
	dev        gopi.SPI
}

const (
	EPD_PIN_RESET = gopi.GPIOPin(17)
	EPD_PIN_CS    = gopi.GPIOPin(8)
	EPD_PIN_DC    = gopi.GPIOPin(25)
	EPD_PIN_BUSY  = gopi.GPIOPin(24)

	EPD_SPI_SPEED = 10000000
	EPD_SPI_BUS   = 0
	EPD_SPI_SLAVE = 0
	EPD_SPI_MODE  = gopi.SPI_MODE_0
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *EPD) Define(cfg gopi.Config) error {
	this.bus = cfg.FlagUint("spi.bus", 0, "SPI Bus")
	this.slave = cfg.FlagUint("spi.slave", 0, "SPI Bus")
	return nil
}

func (this *EPD) New(gopi.Config) error {
	if dev := this.Devices.Open(*bus, *slave); err != nil {
		return err
	} else {
		fmt.Println("spi=", dev)
	}
	return nil
}

func (this *EPD) Dispose() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *EPD) String() string {
	str := "<epd"
	str += " gpio=" + fmt.Sprint(this.GPIO)
	str += " spi=" + fmt.Sprint(this.SPI)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// RESET

func (this *EPD) Reset() error {

}
