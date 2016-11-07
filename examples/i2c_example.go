/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"
	"flag"
)

import (
	gopi "github.com/djthorpe/gopi"
	app "github.com/djthorpe/gopi/app"
	linux "github.com/djthorpe/gopi/device/linux"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// Which I2C controller
	bus := app.FlagSet.Lookup("bus").Value.(flag.Getter).Get().(uint)
	slave := app.FlagSet.Lookup("slave").Value.(flag.Getter).Get().(uint)

	// Create the Pimote interface
	i2c, err := gopi.Open(linux.I2C{ Bus: bus, Slave: uint8(slave) },app.Logger)
	if err != nil {
		return err
	}
	defer i2c.Close()

	app.Logger.Info("I2C=%v",i2c)

	return err
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_I2C)

	// Add on command-line flags
	config.FlagSet.Uint("bus",0,"Bus (0,1 or 2)")
	config.FlagSet.Uint("slave",0,"Slave Address")

	// Create the application
	myapp, err := app.NewApp(config)
	if err == flag.ErrHelp {
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(RunLoop); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
