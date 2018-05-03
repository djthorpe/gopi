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
	hw "github.com/djthorpe/gopi/cmd/rpc/helloworld"
)

////////////////////////////////////////////////////////////////////////////////

func RecordMatches(app *gopi.AppInstance, record *gopi.RPCServiceRecord) bool {
	if service_type, err := gopi.RPCServiceType(app.Service(), 0); err != nil {
		app.Logger.Error("RecordMatches: %v", err)
		return false
	} else {
		return service_type == record.Type
	}
}

func RunClient(app *gopi.AppInstance, client *hw.MyGreeterClient) {
	name, _ := app.AppFlags.GetString("name")
	if message, err := client.SayHello(name); err != nil {
		app.Logger.Error("RunClient: %v", client.Conn().Name(), err)
	} else {
		fmt.Printf("%v says '%v'\n\n", client.Conn().Name(), message)
	}
}

func HandleEvent(app *gopi.AppInstance, evt gopi.RPCEvent) {
	pool := app.ModuleInstance("rpc/clientpool").(gopi.RPCClientPool)

	switch evt.Type() {
	case gopi.RPC_EVENT_SERVICE_RECORD:
		// Create a connection if the service record is correct type
		if RecordMatches(app, evt.ServiceRecord()) {
			if _, err := pool.Connect(evt.ServiceRecord(), 0); err != nil {
				app.Logger.Error("Connect: %v", err)
			}
		}
	case gopi.RPC_EVENT_CLIENT_CONNECTED:
		conn := evt.Source().(gopi.RPCClientConn)
		if client := pool.NewClient("mutablelogic.Helloworld", conn); client == nil {
			app.Logger.Error("Connect: Unable to create client")
		} else {
			RunClient(app, client.(*hw.MyGreeterClient))
			if err := pool.Disconnect(conn); err != nil {
				app.Logger.Error("Disconnect: %v", err)
			}
		}
	case gopi.RPC_EVENT_CLIENT_DISCONNECT:
		// Send a terminate signal to end
		app.SendSignal()
	default:
		app.Logger.Warn("Unhandled event type: %v", evt.Type())
	}
}

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
			if evt != nil {
				HandleEvent(app, evt.(gopi.RPCEvent))
			}
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

	// Name argument
	config.AppFlags.FlagString("name", "", "Your name")

	// Set the RPCServiceRecord for server discovery
	config.Service = "helloworld"

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main, EventLoop))
}
