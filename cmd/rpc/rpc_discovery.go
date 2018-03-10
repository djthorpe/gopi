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
	"errors"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

func MainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	discovery := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery)
	//timeout, _ := app.AppFlags.GetDuration("timeout")
	service, _ := app.AppFlags.GetString("service")

	// Return error if no service
	if service == "" {
		return errors.New("Missing -service parameter (try _smb._tcp)")
	}

	// Return error if no discovery
	if discovery == nil {
		return errors.New("Missing discovery service")
	}

	/*// Discover services on the network
	d := make(chan bool)
	s := make([]*gopi.RPCServiceRecord, 0)
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err := discovery.Browse(ctx, service); err != nil {
		return err
	}

	<-d

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"name", "type", "host", "txt"})
	table.SetAutoFormatHeaders(false)
	table.SetAutoMergeCells(true)
	for _, service := range s {
		table.Append([]string{service.Name, service.Type, fmt.Sprintf("%v:%v", service.Host, service.Port), strings.Join(service.Text, " ")})
	}
	table.Render()
	*/

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the lirc instance
	config := gopi.NewAppConfig("mdns")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, MainLoop))
}
