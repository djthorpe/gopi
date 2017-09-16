/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"errors"
	"fmt"
	"os/user"
)

import (
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

/////////////////////////////////////////////////////////////////////
// main()

func main() {
	var err gopi.Error

	// Open the driver with configuration
	user, _ := user.Current()
	driver, ok := gopi.Open2(MyConfig{Username: user.Username}, nil, &err).(*MyDriver)
	if !ok { // Driver could not be opened
		fmt.Println("Could not open driver, " + err.Error())
		return
	}
	// Close on exit
	defer driver.Close()

	// Perform actions on driver here
	fmt.Println(driver)
}
