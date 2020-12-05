package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

func (this *app) RunVersion(ctx context.Context, stub gopi.PingStub) error {
	version, err := stub.Version(ctx)
	if err != nil {
		return err
	}

	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.Append([]string{
		"Name", version.Name(),
	})
	tag, branch, hash := version.Version()
	if tag != "" {
		table.Append([]string{
			"Tag", tag,
		})
	}
	if branch != "" {
		table.Append([]string{
			"Branch", branch,
		})
	}
	if hash != "" {
		table.Append([]string{
			"Hash", hash,
		})
	}
	table.Append([]string{
		"Go version", version.GoVersion(),
	})
	if t := version.BuildTime(); t.IsZero() == false {
		table.Append([]string{
			"Build time", t.Format(time.RFC3339),
		})
	}

	if services, err := stub.ListServices(ctx); err != nil {
		return err
	} else if len(services) > 0 {
		for _, service := range services {
			table.Append([]string{
				"Services", service,
			})
		}
	}

	table.Render()

	// Return success
	return nil
}

func (this *app) RunPing(ctx context.Context, stub gopi.PingStub) error {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			now := time.Now()
			if err := stub.Ping(ctx); err != nil {
				return err
			} else {
				fmt.Println("ping: ", time.Since(now))
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
