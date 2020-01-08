// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype

import (
	"os"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	ft "github.com/djthorpe/gopi/v2/sys/freetype"
)

type FontManager struct {
}

type fontmanager struct {
	log                 gopi.Logger
	library             ft.FT_Library
	major, minor, patch int
	faces               map[string]gopi.FontFace

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (FontManager) Name() string { return "gopi.fonts.freetype" }

func (config FontManager) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(fontmanager)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// INIT & CLOSE

func (this *fontmanager) Init(config FontManager) error {
	this.Lock()
	defer this.Unlock()

	if library, err := ft.FT_Init(); err != nil {
		return err
	} else {
		this.library = library
		this.major, this.minor, this.patch = ft.FT_Library_Version(library)
	}

	// Return success
	return nil
}

func (this *fontmanager) Close() error {
	this.Lock()
	defer this.Unlock()

	if this.library != ft.FT_Library(nil) {
		if err := ft.FT_Destroy(this.library); err != nil {
			return err
		}
		this.library = ft.FT_Library(nil)
	}

	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION: gopi.FontManager

// Open a font face - first face at index 0 is loaded
func (this *fontmanager) OpenFace(path string) (gopi.FontFace, error) {
	return nil, gopi.ErrNotImplemented
}

// Open a font face - indexed within file of several faces
func (this *fontmanager) OpenFaceAtIndex(path string, index uint) (gopi.FontFace, error) {
	return nil, gopi.ErrNotImplemented

}

// Open font faces at path, checking to see if individual files should
// be opened through a callback function
func (this *fontmanager) OpenFacesAtPath(path string, callback func(manager gopi.FontManager, path string, info os.FileInfo) bool) error {
	return gopi.ErrNotImplemented

}

// Destroy a font face
func (this *fontmanager) DestroyFace(gopi.FontFace) error {
	return gopi.ErrNotImplemented
}

// Return an array of font families which are loaded
func (this *fontmanager) Families() []string {
	return nil
}

// Return open face for filepath
func (this *fontmanager) FaceForPath(path string) gopi.FontFace {
	return nil
}

// Return faces in a family and/or with a particular set of attributes
func (this *fontmanager) Faces(family string, flags gopi.FontFlags) []gopi.FontFace {
	return nil
}
