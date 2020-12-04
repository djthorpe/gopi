package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

func (this *app) RunPing(ctx context.Context) error {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			fmt.Println("ping")
			if err := this.stub.Ping(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}

	// Return success
	return nil
}

func (this *app) RunVersion(ctx context.Context) error {
	version, err := this.stub.Version(ctx)
	if err != nil {
		return err
	}

	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
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
	table.Render()

	// Return success
	return nil
}
