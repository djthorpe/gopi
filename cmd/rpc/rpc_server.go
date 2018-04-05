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

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("service/helloworld:grpc")

	// Set the RPCServiceRecord for server discovery
	config.Service = "helloworld"

	// Run the server and register all the services
	os.Exit(gopi.RPCServerTool(config))
}
