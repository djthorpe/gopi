/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"errors"
	"fmt"
	"os"
)

// CommandLineTool is the basic form of running a command-line
// application, you generally call this from the main() function
func CommandLineTool(config AppConfig, main_task MainTask, background_tasks ...BackgroundTask) int {

	// Create the application
	app, err := NewAppInstance(config)
	if err != nil {
		if err != ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(main_task, background_tasks...); err == ErrHelp {
		config.AppFlags.PrintUsage()
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

// RPCServer runs a set of RPC Services, you generally call this from the main() function
// and ensure to import rpc/server and rpc/discovery modules anonymously into
// your application as well as all your RPC services
func RPCServerTool(config AppConfig, background_tasks ...BackgroundTask) int {
	// Append on "rpc/server" and "rpc/discovery" onto your module configurations
	var err error
	if config.Modules, err = AppendModulesByName(config.Modules, "rpc/server", "rpc/discovery"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}

	// Create the application
	app, err := NewAppInstance(config)
	if err != nil {
		if err != ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application with a main task and background tasks
	if err := app.Run(mainRPCServer, appendTasks(background_tasks, bgRPCServer, bgRPCDiscovery)...); err == ErrHelp {
		config.AppFlags.PrintUsage()
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func appendTasks(tasks []BackgroundTask, append_tasks ...BackgroundTask) []BackgroundTask {
	return append(tasks, append_tasks...)
}

func mainRPCServer(app *AppInstance, done chan<- struct{}) error {

	// Obtain the RPC Server
	server, ok := app.ModuleInstance("rpc/server").(RPCServer)
	if server == nil || ok == false {
		return errors.New("rpc/server missing")
	}

	// Wait for CTRL+C or SIGTERM
	app.Logger.Info("Waiting for CTRL+C or SIGTERM to stop server")
	app.WaitForSignal()

	// Indicate we want to stop the server - shutdown after we have
	// serviced requests
	if err := server.Stop(false); err != nil {
		return err
	}

	// Finish gracefully
	done <- DONE
	return nil
}

func bgRPCServer(app *AppInstance, done <-chan struct{}) error {

	if server, ok := app.ModuleInstance("rpc/server").(RPCServer); server == nil || ok == false {
		return errors.New("rpc/server: missing or invalid")
	} else if modules := ModulesByType(MODULE_TYPE_SERVICE); len(modules) == 0 {
		return errors.New("rpc/server: no RPC services registered")
	} else {

		// Start the server
		if err := server.Start(); err != nil {
			return err
		}

		// wait for done
		<-done
	}

	// Successful completion
	return nil
}

func bgRPCDiscovery(app *AppInstance, done <-chan struct{}) error {

	if server, ok := app.ModuleInstance("rpc/server").(RPCServer); server == nil || ok == false {
		return errors.New("rpc/server: missing or invalid")
	} else if discovery, ok := app.ModuleInstance("rpc/discovery").(RPCServiceDiscovery); discovery == nil || ok == false {
		return errors.New("rpc/discovery: missing or invalid")
	} else {
		// Listen for server started events
		events := server.Subscribe()

	FOR_LOOP:
		for {
			select {
			case evt := <-events:
				if server_event, ok := evt.(RPCEvent); server_event != nil && ok {
					if server_event.Type() == RPC_EVENT_SERVER_STARTED {
						// Register service
						if service := server.Service(app.service); service != nil {
							if err := discovery.Register(service); err != nil {
								app.Logger.Error("rpc/discovery: %v", err)
							}
						}
					}
				}
			case <-done:
				break FOR_LOOP
			}
		}

		// Stop listening for events
		server.Unsubscribe(events)
	}

	return nil
}
