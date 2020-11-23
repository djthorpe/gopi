// +build !ffmpeg

package ffmpeg

import (
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *Manager) New(gopi.Config) error {
	return gopi.ErrNotImplemented
}

func (this *Manager) OpenFile(path string) (gopi.Media, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Manager) Close(gopi.Media) error {
	return gopi.ErrNotImplemented
}
