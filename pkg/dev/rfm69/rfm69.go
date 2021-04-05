package rfm69

import (
	// Modules
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type RFM69 struct {
	gopi.Unit
	gopi.Logger
	gopi.SPI
	gopi.SPIBus

	bus, slave *uint
	version    uint8
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RFM_SPI_MODE      = gopi.SPI_MODE_0
	RFM_SPI_SPEEDHZ   = 8000000 // 4MHz
	RFM_VERSION_VALUE = 0x24
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *RFM69) Define(cfg gopi.Config) error {
	this.bus = cfg.FlagUint("spi.bus", 0, "SPI Bus for RFM69")
	this.slave = cfg.FlagUint("spi.slave", 1, "SPI Slave for RFM69")
	return nil
}

func (this *RFM69) New(gopi.Config) error {
	this.Require(this.SPI, this.Logger)

	// Set SPI
	this.SPIBus = gopi.SPIBus{*this.bus, *this.slave}
	if err := this.SPI.SetMode(this.SPIBus, RFM_SPI_MODE); err != nil {
		return err
	}
	if err := this.SPI.SetMaxSpeedHz(this.SPIBus, RFM_SPI_SPEEDHZ); err != nil {
		return err
	}

	// Get version - and check against expected value
	if version, err := this.GetVersion(); err != nil {
		return gopi.ErrNotFound.WithPrefix("RFM69")
	} else if version != RFM_VERSION_VALUE {
		return gopi.ErrNotFound.WithPrefix("RFM69")
	} else {
		this.version = version
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS
