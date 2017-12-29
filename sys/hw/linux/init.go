// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register FilePoll
	gopi.RegisterModule(gopi.Module{
		Name: "linux/filepoll",
		Type: gopi.MODULE_TYPE_OTHER,
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(FilePoll{}, app.Logger)
		},
	})

	// Register GPIO
	gopi.RegisterModule(gopi.Module{
		Name:     "linux/gpio",
		Requires: []string{"linux/filepoll"},
		Type:     gopi.MODULE_TYPE_GPIO,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagBool("gpio.unexport", true, "Unexport exported pins on exit")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			unexport, _ := app.AppFlags.GetBool("gpio.unexport")
			return gopi.Open(GPIO{
				UnexportOnClose: unexport,
				FilePoll:        app.ModuleInstance("filepoll/linux").(FilePollInterface),
			}, app.Logger)
		},
	})

	// Register I2C
	gopi.RegisterModule(gopi.Module{
		Name: "linux/i2c",
		Type: gopi.MODULE_TYPE_I2C,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("i2c.bus", 1, "I2C Bus")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			bus, _ := app.AppFlags.GetUint("i2c.bus")
			return gopi.Open(I2C{
				Bus: bus,
			}, app.Logger)
		},
	})

	// Register SPI
	gopi.RegisterModule(gopi.Module{
		Name: "linux/spi",
		Type: gopi.MODULE_TYPE_SPI,
		Config: func(config *gopi.AppConfig) {
			config.AppFlags.FlagUint("spi.bus", 0, "SPI Bus")
			config.AppFlags.FlagUint("spi.slave", 0, "SPI Slave")
			config.AppFlags.FlagUint("spi.delay", 0, "SPI Transfer delay in microseconds")
		},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			bus, _ := app.AppFlags.GetUint("spi.bus")
			slave, _ := app.AppFlags.GetUint("spi.slave")
			delay, _ := app.AppFlags.GetUint16("spi.delay")
			return gopi.Open(SPI{
				Bus:   bus,
				Slave: slave,
				Delay: delay,
			}, app.Logger)
		},
	})

}
