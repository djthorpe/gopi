package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/djthorpe/gopi/v3"
)

func (this *app) Remux(ctx context.Context) error {
	count := uint(0)

	// Process files
	if paths, err := GetFileArgs(this.Command.Args()); err != nil {
		return err
	} else if err := this.Walk(ctx, paths, &count, func(path string, info os.FileInfo) error {
		if err := this.RemuxMedia(ctx, path); err != nil && errors.Is(err, context.Canceled) == false {
			if *this.quiet == false {
				this.Logger.Print(filepath.Base(path), ": ", err)
			}
			return nil
		} else {
			return err
		}
	}); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) RemuxMedia(ctx context.Context, path string) error {
	media, err := this.MediaManager.OpenFile(path)
	if err != nil {
		return err
	}
	defer this.MediaManager.Close(media)

	// Only remux if AUDIO or VIDEO
	if media.Flags()&gopi.MEDIA_FLAG_AUDIO == 0 && media.Flags()&gopi.MEDIA_FLAG_VIDEO == 0 {
		return fmt.Errorf("No audio or video contained in media")
	}

	// Determine output path
	out := OutFilename(path, *this.out, nil)
	if filepath.IsAbs(out) == false {
		if absout, err := filepath.Abs(out); err != nil {
			return err
		} else {
			out = absout
		}
	}

	// If input and output paths are the same, return error
	if path == out {
		return fmt.Errorf("Cannot remux to same file")
	}

	// Create output
	dst, err := this.MediaManager.CreateFile(out)
	if err != nil {
		return err
	}
	defer this.MediaManager.Close(dst)

	// Read and write packets (no decoding of data)
	return media.Read(ctx, nil, func(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
		return dst.Write(ctx, packet)
	})
}
