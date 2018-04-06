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
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

var (
	ctx    context.Context
	cancel context.CancelFunc
	start  chan struct{}
)

////////////////////////////////////////////////////////////////////////////////

func serviceType(service, network string) string {
	return "_" + strings.TrimSpace(service) + "._" + network
}

func processEvent(evt gopi.RPCEvent) error {
	if evt.Type() == gopi.RPC_EVENT_SERVICE_RECORD {
		fmt.Println(evt)
	}
	return nil
}

func EventLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	discovery, ok := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery)
	if discovery == nil || ok == false {
		return errors.New("Missing or invalid rpc/discovery module")
	}

	events := discovery.Subscribe()
FOR_LOOP:
	for {
		select {
		case <-done:
			break FOR_LOOP
		case evt := <-events:
			if rpc_evt, ok := evt.(gopi.RPCEvent); rpc_evt != nil && ok {
				if err := processEvent(rpc_evt); err != nil {
					return err
				}
			}
		}
	}

	// Unsubscribe
	discovery.Unsubscribe(events)

	// Return success
	return nil
}

func Browse(app *gopi.AppInstance, done <-chan struct{}) error {

	// Wait for start (the context is created)
	<-start

	// Browse (blocking until timeout or cancel)
	if discovery, ok := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery); discovery == nil || ok == false {
		return errors.New("Missing or invalid rpc/discovery module")
	} else if err := discovery.Browse(ctx, serviceType(app.Service(), "tcp")); err != nil {
		return err
	}

	// Timeout - send SIGTERM
	if proc, err := os.FindProcess(syscall.Getpid()); err != nil {
		return err
	} else if err := proc.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	// Return success
	<-done
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	// Create context
	if timeout, _ := app.AppFlags.GetDuration("timeout"); timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	}

	// Signal browse function to start
	start <- gopi.DONE

	// Wait until CTRL+C is pressed or SIGTERM signal
	app.WaitForSignal()

	// Send cancel to stop the browse function
	cancel()

	// Success
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/discovery")

	// Set the RPCServiceRecord for server discovery
	config.Service = "helloworld"

	// Set start
	start = make(chan struct{})

	// Set flags
	//config.AppFlags.FlagString("addr", "", "Server address")
	config.AppFlags.FlagDuration("timeout", 750*time.Millisecond, "Server discovery timeout")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main, Browse, EventLoop))
}
