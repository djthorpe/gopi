/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Reads bytes from the SPI interface
package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if app.SPI == nil {
		return app.Logger.Error("Missing SPI module instance")
	}

	// Set mode
	if mode, exists := app.AppFlags.GetUint("mode"); exists {
		if err := app.SPI.SetMode(gopi.SPIMode(mode)); err != nil {
			return err
		}
	}

	// Set bits
	if bits, exists := app.AppFlags.GetUint("bits"); exists {
		if err := app.SPI.SetBitsPerWord(uint8(bits)); err != nil {
			return err
		}
	}

	// Set speed
	if speed, exists := app.AppFlags.GetUint("speed"); exists {
		if err := app.SPI.SetMaxSpeedHz(uint32(speed)); err != nil {
			return err
		}
	}

	// Read back values
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Parameter", "Value"})
	table.Append([]string{"mode", fmt.Sprint(app.SPI.Mode())})
	table.Append([]string{"bits_per_word", fmt.Sprint(app.SPI.BitsPerWord())})
	table.Append([]string{"max_speed_hz", fmt.Sprint(app.SPI.MaxSpeedHz())})

	table.Render()
	// Finished
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the spi instance
	config := gopi.NewAppConfig("spi")

	// Flags
	config.AppFlags.FlagUint("mode", 0, "Mode")
	config.AppFlags.FlagUint("speed", 0, "Maximum speed, Hz")
	config.AppFlags.FlagUint("bits", 8, "Bits per word")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
