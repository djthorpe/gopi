/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This sample program shows how you can get information from the GPIO
// device, and set pins input, output, etc. To use the software, there are
// some flags. To enumerate the pins on the device with their current status
// and names:
//
//   gpioctrl
//
// To get status on an individual pin:
//
//   gpioctrl -pin <pin>
//
// The pin can be queried by the physical pin number or by the name of the pin,
// for example GPIO23. To set a pin to OUTPUT and set output high and low:
//
//   gpioctrl -pin <pin> -high
//
//   gpioctrl -pin <pin> -low
//
// To set a pin to INPUT:
//
//   gpioctrl -pin <pin> -input
//
// You can also set a pin mode to be alternate function:
//
//  gpioctrl -pin <pin> -alt 0
//
// And so forth...
//
package main

import (
	"fmt"
	"os"
	"errors"
	"strconv"
)

import (
	app "github.com/djthorpe/gopi/app"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// Check flags
	app.Logger.Info("flags=%v",app.FlagSet)
	err := CheckFlags(app.FlagSet)
	if err != nil {
		return err
	}

	// Debugging output
	app.Logger.Debug("Device=%v", app.Device)
	app.Logger.Debug("GPIO=%v", app.GPIO)

	// Get pin
	gpio := app.GPIO
	pin, err := ParsePinFlag(app.GPIO,app.FlagSet)
	if err != nil {
		return err
	}
	app.Logger.Debug("Pin=%v", pin)

	// If no pin, then print out the table of pin states
	switch {
	case pin == hw.GPIO_PIN_NONE:
		return PrintPinTable(gpio,os.Stdout)
	case app.FlagSet.HasFlag("low"):
		gpio.SetPinMode(pin,hw.GPIO_OUTPUT)
		gpio.WritePin(pin,hw.GPIO_LOW)
		return PrintPinTable(gpio,os.Stdout)
	case app.FlagSet.HasFlag("high"):
		gpio.SetPinMode(pin,hw.GPIO_OUTPUT)
		gpio.WritePin(pin,hw.GPIO_HIGH)
		return PrintPinTable(gpio,os.Stdout)
	case app.FlagSet.HasFlag("input"):
		gpio.SetPinMode(pin,hw.GPIO_INPUT)
		return PrintPinTable(gpio,os.Stdout)
	case app.FlagSet.HasFlag("alt"):
		alt, _ := app.FlagSet.GetUint("alt")
		gpio.SetPinMode(pin,AltMode(alt))
		return PrintPinTable(gpio,os.Stdout)
	default:
		return errors.New("NOT IMPLEMENTED")
	}

	// Return success
	return nil
}

func AltMode(alt uint) hw.GPIOMode {
	switch(alt) {
	case 0:
		return hw.GPIO_ALT0
	case 1:
		return hw.GPIO_ALT1
	case 2:
		return hw.GPIO_ALT2
	case 3:
		return hw.GPIO_ALT3
	case 4:
		return hw.GPIO_ALT4
	case 5:
		return hw.GPIO_ALT5
	default:
		return hw.GPIO_INPUT
	}
}

////////////////////////////////////////////////////////////////////////////////

func PrintPinTable(gpio hw.GPIODriver,fd *os.File) error {

	// print out pin table in two columns
	rows := gpio.NumberOfPhysicalPins() / 2
	header := "+----+----------+----+----------+ +----------+----+----------+----+"
	format := "| %4s | %8s | %6s | %2s | | %-2s | %-6s | %-8s | %-4s |"

	fmt.Fprintln(fd,header)

	for i := uint(0); i < rows; i++ {
		args := make([]interface{},8)
		args[3],args[2],args[1],args[0] = PinStateAsString(gpio,(i * 2) + 1)
		args[4],args[5],args[6],args[7] = PinStateAsString(gpio,(i * 2) + 2)
		fmt.Fprintln(fd,fmt.Sprintf(format,args...))
	}

	fmt.Fprintln(fd,header)
	return nil
}

func PinStateAsString(gpio hw.GPIODriver,physicalpin uint) (string,string,string,string) {
	logicalpin := gpio.PhysicalPin(physicalpin)
	if logicalpin == hw.GPIO_PIN_NONE {
		return strconv.FormatUint(uint64(physicalpin),10),"","",""
	}
	return strconv.FormatUint(uint64(physicalpin),10),logicalpin.String(),gpio.GetPinMode(logicalpin).String(),gpio.ReadPin(logicalpin).String()
}

////////////////////////////////////////////////////////////////////////////////

func CheckFlags(flagset *app.Flags) error {
	// if no flags, then return OK
	if len(flagset.Flags()) == 0 {
		return nil
	}

	// Check for either: low, high, input or alt which are mutually
	// exclusive flags
	c := 0
	for _,flag := range([]string{ "input","alt","low","high" }) {
		if flagset.HasFlag(flag) {
			c++
		}
	}
	if c != 1 {
		return errors.New("One of -low, -high, -input, or -alt required")
	}

	// check for alt being between 0 and 5
	alt, exists := flagset.GetUint("alt")
	if exists {
		if alt > 5 {
			return errors.New("-alt is required to be between 0 and 5")
		}
	}

	// Check for -pin argument
	if flagset.HasFlag("pin") == false {
		return errors.New("-pin flag required")
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func ParsePinFlag(gpio hw.GPIODriver,flagset *app.Flags) (hw.GPIOPin,error) {
	value, exists := flagset.GetString("pin")
	if exists == false {
		return hw.GPIO_PIN_NONE,nil
	}
	pin, err := strconv.ParseUint(value,10,32)
	if err == nil {
		logical := gpio.PhysicalPin(uint(pin))
		if logical == hw.GPIO_PIN_NONE {
			return logical,errors.New("Invalid pin")
		}
		return logical, nil
	}
	return hw.GPIO_PIN_NONE,errors.New("NOT IMPLEMENTED")
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Configuration
	config := app.Config(app.APP_GPIO)

	// Set the flags
	config.FlagSet.FlagString("pin","","Physical Pin Number or name")
	config.FlagSet.FlagBool("low",false,"Set pin to OUTPUT and set pin level LOW")
	config.FlagSet.FlagBool("high",false,"Set pin to OUTPUT and set pin level HIGH")
	config.FlagSet.FlagBool("input",false,"Set pin to INPUT")
	config.FlagSet.FlagUint("alt",0,"Set pin to an alternate function 0-5")

	// Create the application
	myapp, err := app.NewApp(config)
	if err == app.ErrHelp {
		// Help requested
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
