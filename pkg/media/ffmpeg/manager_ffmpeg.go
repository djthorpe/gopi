// +build ffmpeg

package ffmpeg

import (
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Logger
	sync.Mutex
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
		this.Logger.Print(level, message)
	})

	// Initialize format
	ffmpeg.AVFormatInit()

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Deinit
	ffmpeg.AVFormatDeinit()

	// Return to standard logging
	ffmpeg.AVLogSetCallback(0, nil)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - OPEN/CLOSE

func (this *Manager) OpenFile(path string) (gopi.Media, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Manager) Close(gopi.Media) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<mediamanager"
	return str + ">"
}
