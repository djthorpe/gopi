package waveshare

import (
	gopi "github.com/djthorpe/gopi/v3"
	gpiobcm "github.com/djthorpe/gopi/v3/pkg/hw/gpiobcm"
)

type EPD struct {
	gopi.Unit
	*gpiobcm.GPIO
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
