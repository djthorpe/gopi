package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/djthorpe/gopi/v3"
)

func (this *app) Play(ctx context.Context) error {
	count := uint(0)

	// Process files
	if paths, err := GetFileArgs(this.Command.Args()); err != nil {
		return err
	} else if err := this.Walk(ctx, paths, &count, func(path string, info os.FileInfo) error {
		if err := this.PlayMedia(ctx, path); err != nil {
			if *this.quiet == false {
				this.Logger.Print(filepath.Base(path), ": ", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) PlayMedia(ctx context.Context, path string) error {
	media, err := this.MediaManager.OpenFile(path)
	if err != nil {
		return err
	}
	defer this.MediaManager.Close(media)

	// Play first audio stream only
	streams := media.StreamsForFlag(gopi.MEDIA_FLAG_AUDIO)
	if len(streams) == 0 {
		return gopi.ErrBadParameter.WithPrefix("Missing audio stream")
	} else if len(streams) > 1 {
		this.Logger.Debug("There are ", len(streams), " streams but only playing the first one")
	}

	// Iterate through the frames decoding them
	return media.DecodeIterator(ctx, []int{streams[0]}, func(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
		return media.DecodeFrameIterator(ctx, packet, func(frame gopi.MediaFrame) error {
			fmt.Println("f=", frame)
			return nil
		})
	})
}
