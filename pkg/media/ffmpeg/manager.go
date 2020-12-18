// +build ffmpeg

package ffmpeg

import (
	"os"
	"path/filepath"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Logger
	sync.Mutex

	in []*inputctx
	//	out []*outputctx
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func (this *Manager) New(gopi.Config) error {
	if this.Logger == nil {
		return gopi.ErrInternalAppError.WithPrefix("gopi.Logger")
	}
	level := ffmpeg.AV_LOG_ERROR
	if this.Logger.IsDebug() {
		level = ffmpeg.AV_LOG_DEBUG
	}
	ffmpeg.AVLogSetCallback(level, func(level ffmpeg.AVLogLevel, message string, userInfo uintptr) {
		this.Logger.Print(level, " ", message)
	})

	// Initialize format
	ffmpeg.AVFormatInit()

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var result error

	// Close all outputs
	/*for _, out := range this.out {
		if out != nil {
			if err := out.Close(); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}*/

	// Close all inputs
	for _, in := range this.in {
		if in != nil {
			if err := in.Close(); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Deinit
	ffmpeg.AVFormatDeinit()

	// Return to standard logging
	ffmpeg.AVLogSetCallback(0, nil)

	// Release resources
	this.in = nil
	//this.out = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - OPEN/CLOSE

func (this *Manager) OpenFile(path string) (gopi.Media, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Clean up the path
	if filepath.IsAbs(path) == false {
		if path_, err := filepath.Abs(path); err == nil {
			path = filepath.Clean(path_)
		}
	}

	// Create the media object and return it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, gopi.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	if ctx := ffmpeg.NewAVFormatContext(); ctx == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewAVFormatContext")
	} else if err := ctx.OpenInput(path, nil); err != nil {
		// when error is returned free is already called
		return nil, err
	} else if in := NewInputContext(ctx); in == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewInputContext")
	} else {
		this.in = append(this.in, in)
		return in, nil
	}
}

func (this *Manager) CreateFile(path string) (gopi.Media, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	return nil, gopi.ErrNotImplemented
	/*
		// Clean up the path
		if filepath.IsAbs(path) == false {
			if path_, err := filepath.Abs(path); err == nil {
				path = filepath.Clean(path_)
			}
		}

		if ctx, err := ffmpeg.NewAVFormatOutputContext(filename, nil); err != nil {
			return nil, err
		} else if out := NewOutputContext(ctx); out == nil {
			return nil, gopi.ErrInternalAppError.WithPrefix("NewOutputContext")
		} else {
			this.out = append(this.out, out)
			return out, nil
		}*/
}

func (this *Manager) Close(media gopi.Media) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if in, ok := media.(*inputctx); in == nil || ok == false {
		return gopi.ErrInternalAppError.WithPrefix("Close")
	} else {
		for i, in_ := range this.in {
			if in_ != in {
				continue
			}
			err := in.Close()
			this.in[i] = nil
			return err
		}
		return gopi.ErrNotFound.WithPrefix("Close")
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<manager.ffmpeg"
	return str + ">"
}
