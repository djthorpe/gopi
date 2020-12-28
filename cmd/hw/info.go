package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
)

////////////////////////////////////////////////////////////////////////////////

func (this *app) RunInfo(ctx context.Context) error {
	args := this.Command.Args()
	ctx, cancel := context.WithTimeout(ctx, *this.timeout)
	defer cancel()

	if this.Platform == nil {
		return gopi.ErrInternalAppError.WithPrefix("Platform")
	}

	if len(args) != 0 {
		return gopi.ErrHelp
	}
	// Display platform information
	table := table.New(table.WithHeader(false), table.WithMergeCells())

	table.Append(header{"Product"}, this.Platform.Product(), fmt.Sprint(this.Platform.Type()))
	table.Append("Serial Number", "", this.Platform.SerialNumber())
	table.Append("Uptime", "", this.Platform.Uptime().Truncate(time.Second).String())
	if l1, l5, l15 := this.Platform.LoadAverages(); l1 != 0 && l5 != 0 && l15 != 0 {
		table.Append("Load Averages", "1m", fmt.Sprintf("%.2f", l1))
		table.Append("Load Averages", "5m", fmt.Sprintf("%.2f", l5))
		table.Append("Load Averages", "15m", fmt.Sprintf("%.2f", l15))
	}
	if zones := this.Platform.TemperatureZones(); len(zones) > 0 {
		for k, v := range zones {
			table.Append("Temperature Zones", k, fmt.Sprintf("%.2fC", v))
		}
	}
	table.Render(os.Stdout)

	// Return success
	return nil
}
