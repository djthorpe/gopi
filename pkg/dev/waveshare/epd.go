package waveshare

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/gpio/broadcom"
)

type EPD struct {
	gopi.Unit
	gopi.GPIO
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
	return nil
}

func (this *EPD) New(gopi.Config) error {
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
	return str + ">"
}
