/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// The client connects to a remote server
package main

//go:generate protoc helloworld/helloworld.proto --go_out=plugins=grpc:.

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

func MainLoop(app *gopi.AppInstance, done chan struct{}) error {
	if client := app.ModuleInstance("rpc/client").(gopi.RPCClient); client == nil {
		return fmt.Errorf("Missing module: rpc/client")
	} else if err := client.Connect(); err != nil {
		return err
	} else {
		defer client.Disconnect()

		// Do things here
		if modules, err := client.Modules(); err != nil {
			return err
		} else {
			app.Logger.Info("client=%v modules=%v", client, modules)
		}
	}

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {

	// Create the configuration
	config := gopi.NewAppConfig("mdns", "rpc/client")

	// Create the application
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(MainLoop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
