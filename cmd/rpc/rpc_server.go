/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// The server serves the GRPC reflection package and the
// helloworld package, which is described in helloworld/helloworld.proto
// In order to install this package, you will need to run go generate with
// both the protoc compiler and the GRPC GO plugin available:
//
// mac# brew install protobuf
// rpi# sudo apt install protobuf-compiler
// go get -u github.com/golang/protobuf/protoc-gen-go
//
// Then:
//
// cd "${GOPATH}/src/github.com/djthorpe/gopi"
// go generate protobuf
// go install cmd/rpc/rpc_server.go
//
package main

import (
	"errors"
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"

	// RPC Services
	_ "github.com/djthorpe/gopi/cmd/rpc/helloworld"
)

////////////////////////////////////////////////////////////////////////////////

func ServerLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	if server, ok := app.ModuleInstance("rpc/server").(gopi.RPCServer); server == nil || ok == false {
		return errors.New("rpc/server missing")
	} else if modules := gopi.ModulesByType(gopi.MODULE_TYPE_SERVICE); len(modules) == 0 {
		return errors.New("No RPC services")
	} else {
		// Create the application instances
		services := make([]gopi.RPCService, 0, len(modules))
		for _, module := range modules {
			if service, ok := app.ModuleInstance(module.Name).(gopi.RPCService); service == nil || ok == false {
				return fmt.Errorf("Unable to create service: %v", module.Name)
			} else {
				services = append(services, service)
			}
		}
		// Start the server with the services
		if err := server.Start(services...); err != nil {
			return err
		}
		// wait for done
		<-done
	}

	// Bomb out
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the lirc instance
	config := gopi.NewAppConfig()

	// Set the RPC Discovery service name
	config.Service = "helloworld"

	// Run the command line tool
	os.Exit(gopi.RPCServerTool(config, ServerLoop))
}
