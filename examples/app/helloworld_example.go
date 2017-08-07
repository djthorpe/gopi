/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// The canonical hello world example demonstrates printing
// hello world and then exiting. Here we use the 'generic'
// set of modules which provide generic system services
package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/mock"
)

////////////////////////////////////////////////////////////////////////////////

func helloWorld(app *gopi.AppInstance, done chan struct{}) error {
	fmt.Println("Hello, World")
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	if app, err := gopi.NewAppInstance(gopi.NewAppConfig()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	} else if err := app.Run(helloWorld); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}
