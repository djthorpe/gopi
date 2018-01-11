/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Outputs a table of information from "vcgencmd"
package main

import (
	"errors"
	"fmt"
	"os"
	"sort"

	// Frameworks
	"./hw_sys"
	"github.com/djthorpe/gopi"
	"github.com/olekukonko/tablewriter"

	// Modules
	rpi "github.com/djthorpe/gopi/sys/hw/rpi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	if app.Hardware == nil {
		return errors.New("No hardware detected")
	} else if hw, ok := app.Hardware.(rpi.VideoCore); ok == false {
		return errors.New("Raspberry Pi hardware not detected")
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"param", "name", "value"})
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		// output OTP
		if otp, err := hw.GetOTP(); err != nil {
			return err
		} else {
			var keys []int
			for k := range otp {
				keys = append(keys, int(k))
			}
			sort.Ints(keys)
			for _, k := range keys {
				table.Append([]string{"otp", fmt.Sprintf("%02d", k), fmt.Sprintf("0x%08X", otp[byte(k)])})
			}
		}

		// Output Serial number and revision
		if serial_number, err := hw.GetSerialNumberUint64(); err != nil {
			return err
		} else {
			table.Append([]string{"serial_number", "", fmt.Sprintf("0x%X", serial_number)})
		}
		if revision, err := hw.GetRevisionUint32(); err != nil {
			return err
		} else {
			table.Append([]string{"revision", "", fmt.Sprintf("0x%X", revision)})
		}

		// GetCoreTemperatureCelcius gets CPU core temperature in celcius
		if core_temperature, err := hw.GetCoreTemperatureCelcius(); err != nil {
			return err
		} else {
			table.Append([]string{"core_temperature", "", fmt.Sprintf("%.1fC", core_temperature)})
		}

		// render table
		table.Render()
	}

	return nil
}

func main() {
	// Create the configuration, load the gpio instance
	config := gopi.NewAppConfig(hw_sys.MODULE_NAMES)

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
