package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

func (this *app) Metadata(ctx context.Context) error {
	count := uint(0)
	files := []*media{}

	// Process files
	if paths, err := GetFileArgs(this.Command.Args()); err != nil {
		return err
	} else if err := this.Walk(ctx, paths, &count, func(path string, info os.FileInfo) error {
		if media, err := this.ProcessMetadata(path); err != nil {
			if *this.quiet == false {
				this.Logger.Print(filepath.Base(path), ": ", err)
			}
		} else {
			files = append(files, media)
		}
		return nil
	}); err != nil {
		return err
	}

	// Print out stream information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(append([]string{"Name", "Type"}, this.FieldNames()...))
	table.SetAutoFormatHeaders(false)
	for _, file := range files {
		row := []string{filepath.Base(file.URL.Path), FormatFlags(file.Flags)}
		row = append(row, FormatMetadata(file, this.FieldKeys())...)
		table.Append(row)
	}
	table.Render()

	return nil
}

func (this *app) FieldKeys() []gopi.MediaKey {
	// Obtain key/value pairs
	keys := []gopi.MediaKey{}
	for key := range this.fields {
		keys = append(keys, key)
	}
	return keys
}

func (this *app) FieldNames() []string {
	result := []string{}
	for _, key := range this.FieldKeys() {
		result = append(result, this.fields[key])
	}
	return result
}

func (this *app) ProcessMetadata(path string) (*media, error) {
	media, err := this.MediaManager.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer this.MediaManager.Close(media)

	// Create obj
	m := NewMedia(media)

	// Append metadata
	meta := media.Metadata()
	for _, key := range meta.Keys() {
		// Add into list of columns
		if _, exists := this.fields[key]; exists == false {
			this.fields[key] = fmt.Sprint(key)
		}

		// Set metadata value
		m.Meta[key] = meta.Value(key)
	}

	// Return success
	return m, nil
}
