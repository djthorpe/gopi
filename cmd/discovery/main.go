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
	"strings"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
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
	} else if len(services) == 0 {
		return gopi.ErrNotFound
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
	all := []gopi.RPCServiceRecord{}
	for _, service := range services {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if services, err := discovery.Lookup(ctx, service); err != nil {
			return err
		} else {
			all = append(all, services...)
		}
	}
	if len(all) == 0 {
		return gopi.ErrNotFound
	}

	table := tablewriter.NewWriter(os.Stdout)
	for _, service := range all {
		hostPort := fmt.Sprintf("%s:%d", service.Host, service.Port)
		if service.Port == 0 {
			hostPort = ""
		}
		addrs := ""
		for _, addr := range service.Addrs {
			addrs += addr.String() + " "
		}
		table.Append([]string{
			service.Name,
			service.Service,
			hostPort,
			strings.TrimSpace(addrs),
		})
	}
	table.Render()

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

func Watch(app gopi.App) error {
	watch := app.Flags().GetBool("watch", gopi.FLAG_NS_DEFAULT)
	if watch {
		fmt.Println("Press CTRL+C to exit")
		app.WaitForSignal(context.Background(), os.Interrupt)
	}
	return nil
}

func Main(app gopi.App, args []string) error {
	if len(args) == 0 {
		if err := EnumerateServices(app); err != nil {
			return err
		} else if err := Watch(app); err != nil {
			return err
		}
	} else if app.Flags().GetBool("register", gopi.FLAG_NS_DEFAULT) {
		return RegisterServices(app, args)
	} else {
		return LookupServices(app, args)
	}

	// Success
	return nil
}
