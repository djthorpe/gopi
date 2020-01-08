// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype

import (
	"fmt"
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
		this.faces = make(map[string]gopi.FontFace)
	}

	// Return success
	return nil
}

func (this *fontmanager) Close() error {
	this.Lock()
	defer this.Unlock()

	for key, face := range this.faces {
		if err := this.DestroyFace(face); err != nil {
			return err
		} else {
			delete(this.faces, key)
		}
	}

	if this.library != ft.FT_Library(nil) {
		if err := ft.FT_Destroy(this.library); err != nil {
			return err
		}
	}

	// Release resources
	this.library = nil
	this.faces = nil

	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *fontmanager) String() string {
	if this.Unit.Closed {
		return fmt.Sprintf("<%v>", this.Log.Name())
	} else {
		return fmt.Sprintf("<%v version={%v,%v,%v}>", this.Log.Name(), this.major, this.minor, this.patch)
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION: gopi.FontManager

// Open a font face - first face at index 0 is loaded
func (this *fontmanager) OpenFace(path string) (gopi.FontFace, error) {
	return this.OpenFaceAtIndex(path, 0)
}

// Open a font face - indexed within file of several faces
func (this *fontmanager) OpenFaceAtIndex(path string, index uint) (gopi.FontFace, error) {
	this.Lock()
	defer this.Unlock()

	// Create the face
	face := NewFontFaceWithPath(path)

	if handle, err := ft.FT_NewFace(this.library, path, index); err != nil {
		return nil, err
	} else if err := ft.FT_SelectCharmap(handle, ft.FT_ENCODING_UNICODE); err != nil {
		ft.FT_DoneFace(handle)
		return nil, err
	} else {
		face.handle = handle
	}

	// Add face to list of faces
	this.faces[face.Path] = face

	return face, nil
}

// Open font faces at path, checking to see if individual files should
// be opened through a callback function
func (this *fontmanager) OpenFacesAtPath(path string, callback func(manager gopi.FontManager, path string, info os.FileInfo) bool) error {
	return gopi.ErrNotImplemented

}

// Destroy a font face
func (this *fontmanager) DestroyFace(face gopi.FontFace) error {
	this.Lock()
	defer this.Unlock()

	if face_, ok := face.(*fontface); ok == false || face_ == nil {
		return gopi.ErrBadParameter.WithPrefix("face")
	} else if _, exists := this.faces[face_.Path]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("face")
	} else {
		delete(this.faces, face_.Path)
		return ft.FT_DoneFace(face_.handle)
	}
}

// Return an array of font families which are loaded
func (this *fontmanager) Families() []string {
	this.Lock()
	defer this.Unlock()

	families := make(map[string]bool, 0)
	for _, face := range this.faces {
		family := face.Family()
		if _, exists := families[family]; exists {
			continue
		}
		families[family] = true
	}
	familes_ := make([]string, 0, len(families))
	for k := range families {
		familes_ = append(familes_, k)
	}
	return familes_
}

// Return open face for filepath
func (this *fontmanager) FaceForPath(path string) gopi.FontFace {
	return nil
}

// Return faces in a family and/or with a particular set of attributes
func (this *fontmanager) Faces(family string, flags gopi.FontFlags) []gopi.FontFace {
	this.Lock()
	defer this.Unlock()

	faces := make([]gopi.FontFace, 0)
	for _, face := range this.faces {
		if family != "" && family != face.Family() {
			continue
		}
		switch flags {
		case gopi.FONT_FLAGS_STYLE_ANY:
			faces = append(faces, face)
		case gopi.FONT_FLAGS_STYLE_REGULAR, gopi.FONT_FLAGS_STYLE_BOLD, gopi.FONT_FLAGS_STYLE_ITALIC, gopi.FONT_FLAGS_STYLE_BOLDITALIC:
			if face.Flags()&flags == flags {
				faces = append(faces, face)
			}
		}
	}
	return faces
}

