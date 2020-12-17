package main

import (
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
)

func (this *app) PrintVersion(cfg gopi.Config) error {
	version := cfg.Version()
	table := table.New(table.WithHeader(false))
	tag, branch, hash := version.Version()

	table.Append(header{"Name"}, version.Name())
	if tag != "" {
		table.Append(header{"Tag"}, tag)
	}
	if branch != "" {
		table.Append(header{"Branch"}, branch)
	}
	if hash != "" {
		table.Append(header{"Hash"}, hash)
	}
	table.Append(header{"GoVersion"}, version.GoVersion())
	if t := version.BuildTime(); t.IsZero() == false {
		table.Append(header{"BuildTime"}, t.Format(time.RFC3339))
	}
	table.Render(os.Stdout)

	// Return success
	return nil
}
