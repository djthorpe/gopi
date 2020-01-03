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
	"sync"
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

func ServiceRecord(app gopi.App, service string) (gopi.RPCServiceRecord, error) {
	if hostname, err := os.Hostname(); err != nil {
		return gopi.RPCServiceRecord{}, err
	} else {
		return gopi.RPCServiceRecord{
			Name:    service,
			Service: "_gopi._tcp",
			Host:    hostname,
			Port:    8080,
		}, nil
	}
}

func RegisterServices(app gopi.App, services []string) error {
	register := app.UnitInstance("register").(gopi.RPCServiceRegister)
	timeout := app.Flags().GetDuration("timeout", gopi.FLAG_NS_DEFAULT)
	var wait sync.WaitGroup
	cancels := []context.CancelFunc{}
	for _, service := range services {
		ctx, cancel := context.WithCancel(context.Background())
		cancels = append(cancels, cancel)
		if record, err := ServiceRecord(app, service); err != nil {
			return err
		} else {
			wait.Add(1)
			go func(record gopi.RPCServiceRecord) {
				defer wait.Done()
				fmt.Println("Register:", record)
				if err := register.Register(ctx, record); err != nil {
					app.Log().Error(err)
				}
			}(record)
		}
	}

	fmt.Println("Press CTRL+C to exit")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	app.WaitForSignal(ctx, os.Interrupt)

	// Cancel and wait until all goroutines have completed
	for _, cancel := range cancels {
		cancel()
	}
	wait.Wait()

	// Success
	return nil
}

func Main(app gopi.App, args []string) error {
	if len(args) == 0 {
		return EnumerateServices(app)
	} else if app.Flags().GetBool("register", gopi.FLAG_NS_DEFAULT) {
		return RegisterServices(app, args)
	} else {
		return LookupServices(app, args)
	}
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, "discovery", "register"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Flags
		app.Flags().FlagDuration("timeout", time.Second, "Timeout for discovery")
		app.Flags().FlagBool("register", false, "Register service")

		// Run and exit
		os.Exit(app.Run())
	}
}
