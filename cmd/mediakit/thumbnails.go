package main

import (
	"context"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/djthorpe/gopi/v3"
)

func (this *app) Thumbnails(ctx context.Context) error {
	count := uint(0)

	// Process files
	if paths, err := GetFileArgs(this.Command.Args()); err != nil {
		return err
	} else if err := this.Walk(ctx, paths, &count, func(path string, info os.FileInfo) error {
		if err := this.ProcessThumbnails(path); err != nil {
			if *this.quiet == false {
				this.Logger.Print(filepath.Base(path), ": ", err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (this *app) ProcessThumbnails(path string) error {
	media, err := this.MediaManager.OpenFile(path)
	if err != nil {
		return err
	}
	defer this.MediaManager.Close(media)

	// Ignore if only a file
	if media.Flags() == gopi.MEDIA_FLAG_FILE {
		return fmt.Errorf("Not a media file")
	}

	// Get video stream
	streams := media.StreamsForFlag(gopi.MEDIA_FLAG_VIDEO)
	if len(streams) == 0 {
		return fmt.Errorf("No video information found")
	}

	if err := media.DecodeIterator([]int{streams[0]}, func(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
		return media.DecodeFrameIterator(ctx, packet, func(frame gopi.MediaFrame) error {
			num := ctx.Frame()
			return this.ProcessFrame(path, num, frame)
		})
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) ProcessFrame(path string, num int, frame gopi.MediaFrame) error {
	filename := fmt.Sprintf("%06d", num) + ".thumbnail.png"
	out := filepath.Join(os.TempDir(), filename)
	w, err := os.Create(out)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := png.Encode(w, frame); err != nil {
		return err
	} else {
		fmt.Println(frame)
		fmt.Println("  =>", out)
	}
	return nil
}
