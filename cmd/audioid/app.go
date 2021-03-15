package main

import (
	"context"
	"os"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.MediaManager
	gopi.Logger

	filename string
}

func (this *app) Define(cfg gopi.Config) error {
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
	// Open the file, decode the audio
	if file, err := this.OpenFile(this.filename); err != nil {
		return err
	} else if err := this.Decode(ctx, file); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) Decode(ctx context.Context, file gopi.MediaInput) error {
	// Use the first video stream found
	streams := file.StreamsForFlag(gopi.MEDIA_FLAG_AUDIO)
	if len(streams) == 0 {
		return gopi.ErrNotFound.WithPrefix("Audio stream")
	} else {
		this.Print(file.StreamForIndex(streams[0]))
	}

	// Decode frames
	return file.Read(ctx, streams[0:1], func(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
		return file.DecodeFrameIterator(ctx, packet, func(frame gopi.MediaFrame) error {
			this.Print("Decoded", ctx.Stream(), ctx.Frame(), " => ", frame)
			return nil
		})
	})
}
