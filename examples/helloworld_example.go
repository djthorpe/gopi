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
	app "../app"   /* import "github.com/djthorpe/gopi/app" */
)

////////////////////////////////////////////////////////////////////////////////

func HelloWorld(app *app.App) error {
	app.Logger.Info("Hello, World")
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Create the application
	myapp, err := app.NewApp(app.AppConfig{
		Features:  app.APP_DEVICE,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(HelloWorld); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
