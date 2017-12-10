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

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

func MainLoop(app *gopi.AppInstance, done chan struct{}) error {
	if server := app.ModuleInstance("rpc/server").(gopi.RPCServer); server == nil {
		return fmt.Errorf("Missing module: rpc/server")
	} else {

		// Wait for completed
		app.Logger.Debug2("MainLoop: waiting for termination signal")
		app.WaitForSignal()

		// Quit server
		app.Logger.Debug("MainLoop: Stopping RPC server")
		if err := server.Stop(false); err != nil {
			app.Logger.Error("Error: %v", err)
		}
	}

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// ServerLoop starts the RPC Server
func ServerLoop(app *gopi.AppInstance, done chan struct{}) error {
	if server := app.ModuleInstance("rpc/server").(gopi.RPCServer); server == nil {
		return fmt.Errorf("Missing module: rpc/server")
	} else if err := server.Start(); err != nil {
		return err
	}

	app.Logger.Info("ServerLoop: Server stopped")

	// Wait for done signal
	_ = <-done
	return nil
}

////////////////////////////////////////////////////////////////////////////////

// EventLoop receives events from the RPC Server
func EventLoop(app *gopi.AppInstance, done chan struct{}) error {
	if server := app.ModuleInstance("rpc/server").(gopi.RPCServer); server == nil {
		return fmt.Errorf("Missing module: rpc/server")
	} else if discovery := app.ModuleInstance("mdns").(gopi.RPCServiceDiscovery); discovery == nil {
		return fmt.Errorf("Missing module: mdns")
	} else {
		events := server.Events()
		for {
			select {
			case <-done:
				app.Logger.Debug("EventLoop: Received done signal")
				return nil
			case evt := <-events:
				if err := EventProcess(app, server, discovery, evt); err != nil {
					app.Logger.Error("EventLoop: %v: %v", evt.Type(), err)
				}
			}
		}
	}
}

// EventProcess processes events
func EventProcess(app *gopi.AppInstance, server gopi.RPCServer, discovery gopi.RPCServiceDiscovery, evt gopi.RPCEvent) error {
	switch evt.Type() {
	case gopi.RPC_EVENT_SERVER_STARTED:
		// Register service
		name, _ := app.AppFlags.GetString("name")
		if service := server.Service(name); service == nil {
			return fmt.Errorf("Unable to create service record")
		} else if err := discovery.Register(service); err != nil {
			return err
		} else {
			app.Logger.Debug("Registered service: %v", service)
		}
	default:
		app.Logger.Debug("EventLoop: Server event not handled: %v", evt)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the configuration
	config := gopi.NewAppConfig("mdns", "rpc/server")
	// Add the server name
	config.AppFlags.FlagString("name", "RPC Server", "Server name")
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

	// Run the application - the different tasks are:
	// MainLoop waits for the termination signal, ServerLoop starts server and
	// blocks until it is stopped, and EventLoop performs actions based on
	// server events emitted
	if err := app.Run(MainLoop, ServerLoop, EventLoop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
