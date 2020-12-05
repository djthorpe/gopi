package main

import (
	"context"
	"fmt"
	"os"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

func (this *app) RunDiscovery(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, *this.timeout)
	defer cancel()

	args := this.Command.Args()
	if len(args) == 0 {
		return this.RunDiscoveryServices(ctx)
	} else if len(args) == 1 {
		return this.RunDiscoveryLookup(ctx, args[0])
	} else {
		return gopi.ErrHelp
	}
}

func (this *app) RunDiscoveryServices(ctx context.Context) error {
	// Enumerate services
	services, err := this.ServiceDiscovery.EnumerateServices(ctx)
	if err != nil {
		return err
	}

	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.SetHeader([]string{"Service"})
	for _, service := range services {
		table.Append([]string{
			service,
		})
	}
	table.Render()

	// Return success
	return nil
}

func (this *app) RunDiscoveryLookup(ctx context.Context, name string) error {
	// Enumerate services
	records, err := this.ServiceDiscovery.Lookup(ctx, name)
	if err != nil {
		return err
	}

	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.SetHeader([]string{"Service", "Name", "Record"})
	for _, record := range records {
		table.Append([]string{
			record.Service(),
			record.Name(),
			fmt.Sprint(record),
		})
	}
	table.Render()

	// Return success
	return nil
}
