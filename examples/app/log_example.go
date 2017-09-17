/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows you how to create your own logging module
// instead of importing the default logger (which logs to either
// stderr or a file)
package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// Structures

type myConfig struct{}

type myLoggerDriver struct{}

////////////////////////////////////////////////////////////////////////////////
// Code to register and then create a new logger driver

func init() {
	gopi.RegisterModule(gopi.Module{
		Type: gopi.MODULE_TYPE_LOGGER,
		Name: "my_test_logger",
		New:  newLogger,
	})
}

func newLogger(app *gopi.AppInstance) (gopi.Driver, error) {
	fmt.Println("newLogger called which creates a logger...")
	return gopi.Open(myConfig{}, app.Logger)
}

////////////////////////////////////////////////////////////////////////////////
// The logger implementation

func (config myConfig) Open(logger gopi.Logger) (gopi.Driver, error) {
	fmt.Println("myLoggerDriver.Open()")
	return new(myLoggerDriver), nil
}

func (this *myLoggerDriver) Close() error {
	fmt.Println("myLoggerDriver.Close()")
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func helloWorld(app *gopi.AppInstance, done chan struct{}) error {
	fmt.Println("Hello, World")
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig()); err != nil {
		// Check to see if -help has been triggered
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}
	} else if err := app.Run(helloWorld); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
