/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Example command for discovery of RPC microservices using mDNS
package main

import (
	"context"
	"errors"
	"os"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

var (
	lock sync.Mutex
	gctx context.Context
)

func InitContext() {
	lock.Lock()
}

func SetContext(ctx context.Context) {
	gctx = ctx
	defer lock.Unlock()
}

func GetContext() context.Context {
	lock.Lock()
	defer lock.Unlock()
	return gctx
}

////////////////////////////////////////////////////////////////////////////////

func EventLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	discovery := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery)

	// Return error if no discovery
	if discovery == nil {
		return errors.New("Missing discovery service")
	}

	// Subscribe to record discovery
	c := discovery.Subscribe()

FOR_LOOP:
	for {
		select {
		case evt := <-c:
			if rpc_evt, ok := evt.(gopi.RPCEvent); rpc_evt != nil && ok {
				app.Logger.Info("rpc_evt=%v", rpc_evt)
			}
		case <-done:
			break FOR_LOOP
		}
	}

	// Stop listening for events
	discovery.Unsubscribe(c)

	return nil
}

func BrowseLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	if discovery, ok := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery); discovery == nil || ok == false {
		return errors.New("Missing or invalid discovery service")
	} else if err := discovery.Browse(GetContext(), "_smb._tcp"); err != nil {
		return err
	}

	// Wait for done
	_ = <-done
	return nil
}

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	// Set parameters
	timeout, _ := app.AppFlags.GetDuration("timeout")
	ctx, cancel := context.WithCancel(context.Background())
	SetContext(ctx)

	// Wait until CTRL+C is pressed
	app.WaitForSignalOrTimeout(timeout)

	// Perform cancel
	cancel()

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the lirc instance
	config := gopi.NewAppConfig("mdns")

	// Set flags
	config.AppFlags.FlagDuration("timeout", 0, "Browse timeout")

	// Init
	InitContext()

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop, BrowseLoop, EventLoop))
}
