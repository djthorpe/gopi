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
)

import (
	gopi "github.com/djthorpe/gopi"
	app "github.com/djthorpe/gopi/app"
	adafruit "github.com/djthorpe/gopi/device/adafruit"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// Debugging output
	app.Logger.Debug("Device=%v", app.Device)
	app.Logger.Debug("I2C=%v", app.I2C)

	bme280, err := gopi.Open(adafruit.BME280{I2C: app.I2C}, app.Logger)
	if err != nil {
		return err
	}
	defer bme280.Close()

	app.Logger.Debug("bme280=%v", bme280)

	temp, pressure, humidity, err := bme280.(*adafruit.BME280Driver).ReadValues()
	if err != nil {
		return err
	}
	altitude := bme280.(*adafruit.BME280Driver).AltitudeForPressure(pressure,adafruit.BME280_PRESSURE_SEALEVEL)

	fmt.Printf("    TEMP = %.2f C\n",temp)
	fmt.Printf("PRESSURE = %.2f hPa\n",pressure)
	fmt.Printf("HUMIDITY = %.2f %%RH\n",humidity)
	fmt.Printf("ALTITUDE = %.2f m\n",altitude)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_I2C)

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
