package main

import (
	"context"
	"fmt"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.MediaManager
	gopi.Logger

	filename string

	// Flags
	start, end *uint
}

func (this *app) Define(cfg gopi.Config) error {
	this.start = cfg.FlagUint("start", 0, "Start frame")
	this.end = cfg.FlagUint("end", 0, "End frame")
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.MediaManager)

	if args := cfg.Args(); len(args) != 1 {
		return gopi.ErrBadParameter.WithPrefix("Missing filename")
	} else if stat, err := os.Stat(args[0]); err != nil {
		return gopi.ErrBadParameter.WithPrefix(args[0])
	} else if stat.Mode().IsRegular() == false {
		return gopi.ErrBadParameter.WithPrefix(args[0])
	} else {
		this.filename = args[0]
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	// Open the file, decode the video
	if file, err := this.OpenFile(this.filename); err != nil {
		return err
	} else if err := this.Decode(ctx, this.PathTemplate(), file); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) PathTemplate() string {
	path := filepath.Base(this.filename)
	if ext := filepath.Ext(this.filename); ext != "" {
		path = strings.TrimSuffix(path, ext)
	}
	return path + "_%04d"
}

func (this *app) Decode(ctx context.Context, path string, file gopi.MediaInput) error {
	// Use the first video stream found
	streams := file.StreamsForFlag(gopi.MEDIA_FLAG_VIDEO)
	if len(streams) == 0 {
		return gopi.ErrNotFound.WithPrefix("Video stream")
	} else {
		this.Print(file.StreamForIndex(streams[0]))
	}

	// Decode frames
	f := uint(0)
	return file.Read(ctx, streams[0:1], func(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
		return file.DecodeFrameIterator(ctx, packet, func(frame gopi.MediaFrame) error {
			f++
			if *this.start > 0 && f < *this.start {
				return nil
			}
			if *this.end > 0 && f > *this.end {
				// Quit loop
				return io.EOF
			}
			return this.DecodeFrame(fmt.Sprintf(path, f), frame)
		})
	})
}

func (this *app) DecodeFrame(path string, frame gopi.MediaFrame) error {
	// Save frame as a PNG
	path = path + ".png"
	fh, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fh.Close()
	if err := png.Encode(fh, frame); err != nil {
		return err
	}

	this.Printf("Saved frame:", frame, "=>", path)
	return nil
}
