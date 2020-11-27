package main

import (
	"context"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

func (this *app) Streams(ctx context.Context) error {
	count := uint(0)
	files := []*file{}

	// Process files
	if paths, err := GetFileArgs(this.Command.Args()); err != nil {
		return err
	} else if err := Walk(ctx, paths, &count, this.offset, this.limit, func(path string, info os.FileInfo) error {
		if file, err := this.ProcessStreams(path); err != nil {
			this.Logger.Print(file.Name, ": ", err)
		} else {
			files = append(files, file)
		}
		return nil
	}); err != nil {
		return err
	}

	// Print out stream information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Stream"})
	table.SetAutoFormatHeaders(false)
	table.SetAutoMergeCells(true)
	for _, file := range files {
		table.AppendBulk(FormatStreams(file))
	}
	table.Render()

	return nil
}

func (this *app) ProcessStreams(path string) (*file, error) {
	file := NewFile(path)
	media, err := this.MediaManager.OpenFile(path)
	if err != nil {
		return file, err
	}
	defer this.MediaManager.Close(media)

	// Append streams
	file.Streams = media.Streams()

	// Return success
	return file, nil
}

func FormatStreams(f *file) [][]string {
	rows := [][]string{}
	for _, stream := range f.Streams {
		rows = append(rows, []string{f.Name, fmt.Sprint(stream)})
	}
	return rows
}
