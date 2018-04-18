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
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi"
	tablewriter "github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/rpc"
)

////////////////////////////////////////////////////////////////////////////////

func MainLoop(app *gopi.AppInstance, done chan struct{}) error {
	mdns := app.ModuleInstance("rpc/discovery").(gopi.RPCServiceDiscovery)
	timeout, _ := app.AppFlags.GetDuration("timeout")
	service, _ := app.AppFlags.GetString("service")

	// Return error if no service
	if service == "" {
		return errors.New("Missing -service parameter (try _smb._tcp)")
	}

	// Discover services on the network
	d := make(chan bool)
	s := make([]*gopi.RPCService, 0)
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err := mdns.Browse(ctx, service, func(service *gopi.RPCService) {
		if service != nil {
			s = append(s, service)
		} else {
			d <- true
		}
	}); err != nil {
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

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

func registerFlags(config gopi.AppConfig) gopi.AppConfig {
	// Register the -name flag
	config.AppFlags.FlagString("service", "", "Service type")
	config.AppFlags.FlagDuration("timeout", time.Second*2, "Timeout")
	// Return config
	return config
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the application
	app, err := gopi.NewAppInstance(registerFlags(gopi.NewAppConfig("mdns")))
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(MainLoop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
