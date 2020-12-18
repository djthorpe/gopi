package main

import (
	"context"
	"image"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djthorpe/gopi/v3"
	"golang.org/x/image/draw"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "github.com/djthorpe/gopi/v3/pkg/dev/waveshare"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/spi"
)

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.EPD

	args     []string
	loop     *bool
	interval *time.Duration
}

func (this *app) Define(cfg gopi.Config) error {
	this.loop = cfg.FlagBool("loop", false, "Loop display")
	this.interval = cfg.FlagDuration("interval", 60*time.Second, "Display change interval")
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	this.args = cfg.Args()
	return nil
}

func (this *app) Run(ctx context.Context) error {
	if this.EPD == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing EPD interface")
	}
	if len(this.args) == 0 {
		return this.EPD.Clear(ctx)
	}

	ticker := time.NewTicker(*this.interval)
	timer := time.NewTimer(time.Millisecond)
	defer ticker.Stop()
	defer timer.Stop()
	i := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if err := this.Draw(ctx, this.args[i]); err != nil {
				return err
			}
			i++
			if i == len(this.args) {
				if *this.loop {
					i = 0
				} else {
					return nil
				}
			}
		case <-ticker.C:
			if err := this.Draw(ctx, this.args[i]); err != nil {
				return err
			}
			i++
			if i == len(this.args) {
				if *this.loop {
					i = 0
				} else {
					return nil
				}
			}
		}
	}
}

func (this *app) Draw(ctx context.Context, path string) error {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".gif", ".png", ".jpg", ".jpeg":
		if fh, err := os.Open(path); err != nil {
			return err
		} else {
			defer fh.Close()
			this.Print(filepath.Base(path))
			if image, _, err := image.Decode(fh); err != nil {
				return err
			} else if err := this.DrawImage(ctx, image); err != nil {
				return err
			}
		}
	default:
		return gopi.ErrBadParameter.WithPrefix(filepath.Base(path))
	}

	// Return success
	return nil
}

func (this *app) DrawImage(ctx context.Context, src image.Image) error {
	size := this.EPD.Size()
	dst := image.NewRGBA(image.Rect(0, 0, int(size.W), int(size.H)))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return this.EPD.Draw(ctx, dst)
}
