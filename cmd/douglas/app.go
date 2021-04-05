package main

import (
	"context"
	"image"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.EPD

	args     []string
	files    []string
	loop     *bool
	interval *time.Duration
	scale    *float64
}

func (this *app) Define(cfg gopi.Config) error {
	this.loop = cfg.FlagBool("loop", false, "Loop display")
	this.interval = cfg.FlagDuration("interval", 60*time.Second, "Display change interval")
	this.scale = cfg.FlagFloat("scale", 1.0, "Scale image in frame")
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

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if path, err := this.Cycle(); err != nil {
				return err
			} else if err := this.Draw(ctx, path); err != nil {
				return err
			}
		case <-ticker.C:
			if path, err := this.Cycle(); err != nil {
				return err
			} else if err := this.Draw(ctx, path); err != nil {
				return err
			}
		}
	}
}

func (this *app) Cycle() (string, error) {
	var path string

	// Populate files array
	if len(this.files) == 0 {
		this.files = this.filewalk(this.args)
	}
	if len(this.files) == 0 {
		return path, context.Canceled
	}

	// Pop first file
	path, this.files = this.files[0], this.files[1:]

	// Check for looping condition
	if len(this.files) == 0 && *this.loop == false {
		return path, context.Canceled
	} else {
		return path, nil
	}
}

func (this *app) filewalk(patterns []string) []string {
	paths := []string{}
	for _, path := range patterns {
		matches, err := filepath.Glob(path)
		if err != nil {
			this.Printf("%q: %v", path, err)
			continue
		} else if len(matches) == 0 {
			this.Printf("%q: %v", path, gopi.ErrNotFound)
			continue
		}
		for _, match := range matches {
			if strings.HasPrefix(match, ".") {
				continue
			}
			if stat, err := os.Stat(match); err != nil {
				this.Printf("%q: %v", path, err)
				continue
			} else if stat.Mode().IsRegular() == false {
				continue
			} else if evaluateFile(match) == false {
				continue
			} else {
				paths = append(paths, match)
			}
		}
	}
	return paths
}

func evaluateFile(path string) bool {
	fh, err := os.Open(path)
	if err != nil {
		return false
	}
	defer fh.Close()
	data := make([]byte, 512)
	if _, err := fh.Read(data); err != nil {
		return false
	} else {
		mimetype := http.DetectContentType(data)
		return strings.HasPrefix(mimetype, "image/")
	}
}

func (this *app) Draw(ctx context.Context, path string) error {
	fh, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fh.Close()

	this.Print("Display: ", filepath.Base(path))
	if image, _, err := image.Decode(fh); err != nil {
		return err
	} else if err := this.DrawImage(ctx, image); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) DrawImage(ctx context.Context, src image.Image) error {
	return this.EPD.DrawSized(ctx, *this.scale, src)
}
