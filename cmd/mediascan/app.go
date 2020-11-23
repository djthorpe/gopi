package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

type app struct {
	gopi.Unit
	gopi.MediaManager
	gopi.Logger

	folder string
	files  []*file
}

type file struct {
	Name     string
	Error    error
	Flags    string
	Metadata []string
	Streams  []string
}

func (this *app) New(cfg gopi.Config) error {
	if args := cfg.Args(); len(args) != 1 {
		return gopi.ErrHelp
	} else {
		this.folder = args[0]
	}
	return nil
}

func (this *app) Run(ctx context.Context) error {
	_, err := os.Stat(this.folder)
	if err != nil {
		return err
	}

	this.files = []*file{}
	if err := filepath.Walk(this.folder, func(path string, info os.FileInfo, err error) error {
		return this.WalkFunc(ctx, path, info, err)
	}); err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	for _, file := range this.files {
		if file.Error != nil {
			table.Append([]string{
				strconv.Quote(file.Name),
				fmt.Sprint(file.Error),
			})
		} else {
			table.Append([]string{
				strconv.Quote(file.Name),
				file.Flags,
				strings.Join(file.Metadata, "\n"),
				strings.Join(file.Streams, "\n"),
			})
		}
	}
	table.Render()

	// Return success
	return nil
}

func (this *app) WalkFunc(ctx context.Context, path string, info os.FileInfo, err error) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		this.Print("Error", err)
		return err
	}
	if strings.HasPrefix(info.Name(), ".") {
		return nil
	}
	if info.Mode().IsRegular() {
		this.files = append(this.files, this.Read(path, info))
	}
	return nil
}

func (this *app) Read(path string, info os.FileInfo) *file {
	f := &file{Name: info.Name()}
	if media, err := this.MediaManager.OpenFile(path); err != nil {
		this.Print("Error", err)
		f.Error = err
	} else {
		defer this.MediaManager.Close(media)
		flags := strings.Split(fmt.Sprint(media.Flags()), "|")
		f.Flags = strings.Join(flags, ", ")
		for _, k := range media.Metadata().Keys() {
			if k == "id3v2_priv.www.amazon.com" || k == "iTunSMPB" || k == "iTunNORM" || k == "iTunes_CDDB_1" {
				continue
			}
			v := media.Metadata().Value(k)
			v_ := fmt.Sprint(v)
			if _, ok := v.(string); ok {
				v_ = strconv.Quote(v_)
			}
			f.Metadata = append(f.Metadata, fmt.Sprint(k, "=", v_))
		}
		for _, s := range media.Streams() {
			f.Streams = append(f.Streams, fmt.Sprint(s))
		}
	}
	return f
}
