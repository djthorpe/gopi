package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/djthorpe/gopi/v3/pkg/table"
)

func (this *app) Streams(ctx context.Context) error {
	count := uint(0)
	files := []*media{}

	// Process files
	if paths, err := GetFileArgs(this.Command.Args()); err != nil {
		return err
	} else if err := this.Walk(ctx, paths, &count, func(path string, info os.FileInfo) error {
		if media, err := this.ProcessStreams(path, info); err != nil {
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
	t := table.New()
	for _, file := range files {
		t.Add(file.Dict())
	}
	if *this.csv {
		t.RenderCSV(os.Stdout)
	} else {
		t.Render(os.Stdout, table.WithFooter(true))
	}

	// Return success
	return nil
}

func (this *app) ProcessStreams(path string, info os.FileInfo) (*media, error) {
	media, err := this.MediaManager.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer this.MediaManager.Close(media)

	// Create obj
	m := NewMedia(media, info)

	// Append streams
	for _, stream := range media.Streams() {
		m.Streams = append(m.Streams, NewStream(stream))
	}

	// Return success
	return m, nil
}
