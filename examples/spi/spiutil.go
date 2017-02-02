/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// SPIUTIL
//
// This example demonstrates reading from an SPI device
package main

import (
	"fmt"
	"os"
)

import (
	app "github.com/djthorpe/gopi/app"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {
	// Debugging output
	app.Logger.Debug("SPI=%v", app.SPI)

	// Set mode
	if mode, exists := app.FlagSet.GetUint("mode"); exists {
		if err := app.SPI.SetMode(hw.SPIMode(mode)); err != nil {
			return err
		}
	}

	// Set bits
	if bits, exists := app.FlagSet.GetUint("bits"); exists {
		if err := app.SPI.SetBitsPerWord(uint8(bits)); err != nil {
			return err
		}
	}

	// Set speed
	if speed, exists := app.FlagSet.GetUint("speed"); exists {
		if err := app.SPI.SetMaxSpeedHz(uint32(speed)); err != nil {
			return err
		}
	}

	// Print information
	fmt.Printf("Mode = %v\n",app.SPI.GetMode())
	fmt.Printf("Bits per word = %v\n",app.SPI.GetBitsPerWord())
	fmt.Printf("Speed = %vHz\n",app.SPI.GetMaxSpeedHz())

	bytes, err := app.SPI.Transfer([]byte{ 0x00 })
	if err != nil {
		return err
	}
	fmt.Printf("Bytes = %v\n",bytes)

	return nil
}

func main() {
	// Create the config
	config := app.Config(app.APP_SPI)

	// Flags
	config.FlagSet.FlagUint("mode",0, "Mode")
	config.FlagSet.FlagUint("speed",0, "Maximum speed")
	config.FlagSet.FlagUint("bits", 8, "Bits per word")

	// Create the application
	myapp, err := app.NewApp(config)
	if err == app.ErrHelp {
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
