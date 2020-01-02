/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	tablewriter "github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func EnumerateServices(app gopi.App) error {
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	timeout := app.Flags().GetDuration("timeout", gopi.FLAG_NS_DEFAULT)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if services, err := discovery.EnumerateServices(ctx); err != nil {
		return err
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		for _, service := range services {
			table.Append([]string{service})
		}
		table.Render()
	}

	// Return success
	return nil
}

func LookupServices(app gopi.App, services []string) error {
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	timeout := app.Flags().GetDuration("timeout", gopi.FLAG_NS_DEFAULT)
	for _, service := range services {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if services, err := discovery.Lookup(ctx, service); err != nil {
			return err
		} else {
			table := tablewriter.NewWriter(os.Stdout)
			for _, service := range services {
				table.Append([]string{
					service.Name,
					service.Service,
					fmt.Sprintf("%s:%d", service.Host, service.Port),
				})
			}
			table.Render()
		}
	}

	// Return success
	return nil
}

func Main(app gopi.App, args []string) error {
	if len(args) == 0 {
		return EnumerateServices(app)
	} else {
		return LookupServices(app, args)
	}
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, "discovery"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Timeout flag
		app.Flags().FlagDuration("timeout", time.Second, "Timeout for discovery")
		// Run and exit
		os.Exit(app.Run())
	}
}
