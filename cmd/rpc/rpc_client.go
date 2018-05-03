/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// An RPC Server tool, import the services as modules
package main

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc/grpc"
	_ "github.com/djthorpe/gopi/sys/rpc/mdns"

	// RPC Clients
	_ "github.com/djthorpe/gopi/cmd/rpc/helloworld"
)

////////////////////////////////////////////////////////////////////////////////

func EventLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	// Obtain client connection
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	if pool == nil {
		return gopi.ErrAppError
	}

	// Subscribe to events
	poolevents := pool.Subscribe()

FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case evt := <-poolevents:
			fmt.Println("EVENT=", evt)
		}
	}

	// Unsubscribe from events
	pool.Unsubscribe(poolevents)

	// Return success
	return nil
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Wait until CTRL+C is pressed or SIGTERM signal
	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()

	// Success
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/client/helloworld:grpc")

	// Set the RPCServiceRecord for server discovery
	config.Service = "helloworld"

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main, EventLoop))
}
