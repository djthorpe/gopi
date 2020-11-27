// +build !freetype

package freetype

import (
	"os"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type FontManager struct {
	gopi.Unit
}

////////////////////////////////////////////////////////////////////////////////
// INIT & CLOSE

func (this *FontManager) New(gopi.Config) error {
	return nil
}

func (this *FontManager) Dispose() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *FontManager) String() string {
	return "<fontmanager>"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION: gopi.FontManager

func (this *FontManager) OpenFace(path string) (gopi.FontFace, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *FontManager) OpenFaceAtIndex(path string, index uint) (gopi.FontFace, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *FontManager) OpenFacesAtPath(path string, callback func(manager gopi.FontManager, path string, info os.FileInfo) bool) error {
	return gopi.ErrNotImplemented
}

func (this *FontManager) DestroyFace(face gopi.FontFace) error {
	return gopi.ErrNotImplemented
}

func (this *FontManager) Families() []string {
	return nil
}

func (this *FontManager) Faces(family string, flags gopi.FontFlags) []gopi.FontFace {
	return nil
}
