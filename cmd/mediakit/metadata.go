package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/djthorpe/gopi/v3/pkg/table"
)

/////////////////////////////////////////////////////////////////////

func (this *app) Metadata(ctx context.Context) error {
	count := uint(0)
	files := []*media{}

	// Process files
	if paths, err := GetFileArgs(this.Command.Args()); err != nil {
		return err
	} else if err := this.Walk(ctx, paths, &count, func(path string, info os.FileInfo) error {
		if media, err := this.ProcessMetadata(path, info); err != nil {
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

	return nil
}

func (this *app) ProcessMetadata(path string, info os.FileInfo) (*media, error) {
	media, err := this.MediaManager.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer this.MediaManager.Close(media)

	// Create obj
	file := NewMedia(media, info)

	// Append metadata
	meta := media.Metadata()
	for _, key := range meta.Keys() {
		file.Meta[key] = meta.Value(key)
	}

	// Return success
	return file, nil
}
