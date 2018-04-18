/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// I2CREAD
//
// This example demonstrates reading from an I2C device. You can use it
// like this:
//
// i2cread -i2cbus <bus> -slave <addr>
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

import (
	app "github.com/djthorpe/gopi/app"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// Debugging output
	app.Logger.Debug("Device=%v", app.Device)
	app.Logger.Debug("I2C=%v", app.I2C)

	addr, err := GetSlaveAddress(app.FlagSet)
	if err != nil {
		return err
	}

	reg, err := GetRegister(app.FlagSet)
	if err != nil {
		return err
	}

	// Detect slave
	detected, err := app.I2C.DetectSlave(addr)
	if err != nil {
		return err
	}
	if detected == false {
		return errors.New("No device detected at that slave address")
	}

	// Set slave
	if err := app.I2C.SetSlave(addr); err != nil {
		return err
	}

	// Read slave
	byte, err := app.I2C.ReadByte(reg)
	if err != nil {
		return err
	}

	app.Logger.Info("Slave=%02X", addr)
	app.Logger.Info("Byte=%02X", byte)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func GetSlaveAddress(flags *app.Flags) (uint8, error) {
	value, exists := flags.GetString("slave")
	if exists == false {
		return uint8(0), errors.New("Missing -slave parameter")
	}
	addr, err := strconv.ParseUint("0x"+value, 0, 64)
	if err != nil || addr > 0x7F {
		return uint8(0), errors.New("Invalid -slave parameter")
	}
	return uint8(addr), nil
}

func GetRegister(flags *app.Flags) (uint8, error) {
	value, exists := flags.GetString("reg")
	if exists == false {
		return uint8(0), errors.New("Missing -reg parameter")
	}
	reg, err := strconv.ParseUint("0x"+value, 0, 64)
	if err != nil || reg > 0xFF {
		return uint8(0), errors.New("Invalid -reg parameter")
	}
	return uint8(reg), nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_I2C)

	// Flags
	config.FlagSet.FlagString("slave", "", "Slave address")
	config.FlagSet.FlagString("reg", "", "Slave register")

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
