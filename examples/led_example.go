/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This sample program shows how you can use the LED device to blink an array
// of LEDs
package main

import (
	"fmt"
	"os"
	"time"
)

import (
	gopi "github.com/djthorpe/gopi"
	app "github.com/djthorpe/gopi/app"
	hw "github.com/djthorpe/gopi/hw"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// Create the LED array
	pins := []hw.GPIOPin{ app.GPIO.PhysicalPin(40),app.GPIO.PhysicalPin(38),app.GPIO.PhysicalPin(37),app.GPIO.PhysicalPin(36) }
	led, err := gopi.Open(hw.LED{GPIO: app.GPIO, Pins: pins}, app.Logger)
	if err != nil {
		return err
	}
	defer led.Close()

	app.Logger.Info("LED=%v", led)

	go func() {
		// Blink 100ms on/50ms off
		for {
			app.Logger.Debug("ON")
			led.(hw.LEDDriver).On()
			time.Sleep(100 * time.Millisecond)
			app.Logger.Debug("OFF")
			led.(hw.LEDDriver).Off()
			time.Sleep(50 * time.Millisecond)
		}
	}()

	app.WaitUntilDone()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the application
	myapp, err := app.NewApp(app.AppConfig{
		Features: app.APP_GPIO,
		LogLevel: util.LOG_ANY,
	})
	if err != nil {
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
