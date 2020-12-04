package main

import (
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

func (this *app) PrintVersion(cfg gopi.Config) error {
	version := cfg.Version()
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
