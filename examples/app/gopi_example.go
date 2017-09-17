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

////////////////////////////////////////////////////////////////////////////////

func main() {
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else if err := app.Run(mainTask); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

/////////////////////////////////////////////////////////////////////
// mainTask

func mainTask(app *gopi.AppInstance, done chan struct{}) error {
	var err gopi.Error

	// Open the driver with configuration
	user, _ := user.Current()
	driver, ok := gopi.Open2(MyConfig{Username: user.Username}, nil, &err).(*MyDriver)
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
