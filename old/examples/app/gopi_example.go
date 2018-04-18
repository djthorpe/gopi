/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows you how to create a gopi driver
// package, which is a concrete instance of usually an
// abstract interface
package main

import (
	"errors"
	"fmt"
	"os"
	"os/user"

	gopi "github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

/////////////////////////////////////////////////////////////////////
// Configuration and Driver

type MyConfig struct {
	Username string
	Error    bool
}

type MyDriver struct {
	username string
	log      gopi.Logger
}

/////////////////////////////////////////////////////////////////////
// Implement Open and Close methods

func (config MyConfig) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Info("Open")

	driver := new(MyDriver)
	driver.log = log
	driver.username = config.Username

	if config.Error {
		return nil, errors.New("An error occurred")
	}

	return driver, nil
}

func (this *MyDriver) Close() error {
	this.log.Info("Close")
	return nil
}

/////////////////////////////////////////////////////////////////////
// Stringify

func (this *MyDriver) String() string {
	return fmt.Sprintf("MyDriver{ username=%v }", this.username)
}

/////////////////////////////////////////////////////////////////////
// mainTask

func mainTask(app *gopi.AppInstance, done chan struct{}) error {
	var err gopi.Error

	// Open the driver with configuration
	user, _ := user.Current()
	driver, ok := gopi.Open2(MyConfig{Username: user.Username}, app.Logger, &err).(*MyDriver)
	if !ok {
		return fmt.Errorf("Could not open driver: %v", err)
	}
	// Close on exit
	defer driver.Close()

	// Perform actions on driver here
	fmt.Println(driver)

	// return success
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the application
	app, err := gopi.NewAppInstance(gopi.NewAppConfig())
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(mainTask); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
