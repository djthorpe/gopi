package main

import (
	"context"
	"fmt"
	"os"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

func (this *app) Codecs(ctx context.Context) error {
	name := ""
	args := this.Command.Args()
	if len(args) == 1 {
		name = args[0]
	} else if len(args) > 1 {
		return gopi.ErrHelp
	}

	codecs := this.MediaManager.ListCodecs(name, 0)
	if len(codecs) == 0 {
		return gopi.ErrNotFound
	}

	// Print out stream information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "Flags"})
	table.SetAutoFormatHeaders(false)

	for _, codec := range codecs {
		fmt.Println(codec)
		table.Append([]string{
			codec.Name(),
			codec.Description(),
			FormatFlags(codec.Flags()),
		})
	}

	table.Render()
	return nil
}
