package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
)

func (this *app) RunVersion(ctx context.Context, stub gopi.PingStub) error {
	version, err := stub.Version(ctx)
	if err != nil {
		return err
	}

	// Display platform information
	table := table.New(table.WithHeader(false), table.WithMergeCells())

	table.Append(header{"Name"}, version.Name())

	tag, branch, hash := version.Version()
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

	if services, err := stub.ListServices(ctx); err == nil {
		for _, service := range services {
			table.Append(header{"Services"}, service)
		}
	}
	table.Render(os.Stdout)

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
