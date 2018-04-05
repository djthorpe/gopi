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
	"errors"
	"os"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

func BrowseLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	/* TODO: Only start browsing if addr is not set */
	if discovery, ok := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery); discovery == nil || ok == false {
		return errors.New("Missing or invalid rpc/discovery module")
	} else if err := discovery.Browse(GetContext(), "_helloworld._tcp"); err != nil {
		return err
	}

	// Wait for done
	_ = <-done
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Wait until CTRL+C is pressed or SIGTERM signal
	app.WaitForSignal()

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/discovery")

	// Set the RPCServiceRecord for server discovery
	config.Service = "helloworld"

	// Set flags
	config.AppFlags.FlagDuration("addr", "", "Server address")
	config.AppFlags.FlagDuration("timeout", 5*time.Second, "Server discovery timeout")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main, BrowseLoop))
}
