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
	_ "github.com/djthorpe/gopi/sys/rpc"
	_ "github.com/djthorpe/gopi/sys/rpc/mdns"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Obtain client connection
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)
	if conn, err := pool.Connect(0); err != nil {
		return fmt.Errorf("pool.Connect: %v", err)
	} else {
		fmt.Println("We have a conn object, addr = %v", conn.Addr())
	}

	// Wait until CTRL+C is pressed or SIGTERM signal
	app.WaitForSignal()

	// Success
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/clientpool")

	// Set the RPCServiceRecord for server discovery
	config.Service = "helloworld"

	// Set flags
	//config.AppFlags.FlagString("addr", "", "Server address")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
