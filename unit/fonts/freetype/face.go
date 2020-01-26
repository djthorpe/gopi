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
	"image"
	"path"
	"path/filepath"
	"strconv"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	ft "github.com/djthorpe/gopi/v2/sys/freetype"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type fontface struct {
	Path   string
	handle ft.FT_Face
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewFontFaceWithPath(path string) *fontface {
	return &fontface{Path: filepath.Clean(path)}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.FontFace

// Get Face Name (from the filename)
func (this *fontface) Name() string {
	return path.Base(this.Path)
}

// Get Face Index
func (this *fontface) Family() string {
	return ft.FT_FaceFamily(this.handle)
}

func (this *fontface) Style() string {
	return ft.FT_FaceStyle(this.handle)
}

func (this *fontface) Index() uint {
	return ft.FT_FaceIndex(this.handle)
}

// Get Number of faces within the file
func (this *fontface) NumFaces() uint {
	return ft.FT_FaceNumFaces(this.handle)
}

// Number of glyphs for the face
func (this *fontface) NumGlyphs() uint {
	return ft.FT_FaceNumGlyphs(this.handle)
}

// Return properties for face
func (this *fontface) Flags() gopi.FontFlags {
	// If neither bold nor italic, set regular flag
	if flags := ft.FT_FaceStyleFlags(this.handle); flags&gopi.FONT_FLAGS_STYLE_BOLDITALIC == 0 {
		return gopi.FONT_FLAGS_STYLE_REGULAR
	} else {
		return flags
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *fontface) String() string {
	if this.handle == nil {
		return "<gopi.fonts.face name=" + strconv.Quote(this.Name()) + ">"
	} else {
		return "<gopi.fonts.face name=" + strconv.Quote(this.Name()) + " index=" + fmt.Sprint(this.Index()) +
			" family=" + strconv.Quote(this.Family()) +
			" style=" + strconv.Quote(this.Style()) +
			" flags=" + fmt.Sprint(this.Flags()) +
			" num_faces=" + fmt.Sprint(this.NumFaces()) +
			" num_glyphs=" + fmt.Sprint(this.NumGlyphs()) +
			">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// RUNE BITMAPS
	
// Return a bitmap for a rune at pixel size
func (this *fontface) BitmapForRunePixels(ch rune,pixels uint) (image.Image, error) {
	// Set charset as UTF-8
	if err := ft.FT_SelectCharmap(this.handle, ft.FT_ENCODING_UNICODE); err != nil {
		return nil, err
	}
	// Set pixel size
	if err := ft.FT_SetPixelSizes(this.handle, pixels); err != nil {
		return nil, err
	}
	// Load glyph
	if bitmap, x, y, err := ft.FT_Load_Glyph(this.handle,ch, ft.FT_RENDER_MODE_NORMAL); err != nil {
		return nil, err
	} else {
		return NewBitmap(bitmap, x, y)
	}
}

// Return a bitmap for a rune at pixel size (with pixels per inch)
func (this *fontface) BitmapForRunePoints(ch rune,points float32,ppi uint) (image.Image, error) {
	// Set charset as UTF-8
	if err := ft.FT_SelectCharmap(this.handle, ft.FT_ENCODING_UNICODE); err != nil {
		return nil, err
	}
	// Set point size
	if err := ft.FT_SetCharSize(this.handle, points,ppi); err != nil {
		return nil, err
	}
	// Load glyph
	if bitmap, x, y, err := ft.FT_Load_Glyph(this.handle,ch, ft.FT_RENDER_MODE_NORMAL); err != nil {
		return nil, err
	} else {
		return NewBitmap(bitmap, x, y)
	}	
}
