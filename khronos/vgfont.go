/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// VGFONT
//
// This package defines abstract fonts which can be loaded from file
// and rendered on an EGL surface (either bitmap or vector)
//
package khronos /* import "github.com/djthorpe/gopi/khronos" */

import (
	"os"
)

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract font interface
type VGFontDriver interface {
	// Inherit general driver interface
	gopi.Driver

	// Open a font face - first face at index 0 is loaded
	OpenFace(path string) (VGFace, error)

	// Open a font face - indexed within file of several faces
	OpenFaceAtIndex(path string, index uint) (VGFace, error)

	// Open font faces at path, checking to see if individual files should
	// be opened through a callback function
	OpenFacesAtPath(path string, callback func(path string, info os.FileInfo) bool) error

	// Destroy a font face
	DestroyFace(VGFace) error

	// Return an array of font families which are loaded
	GetFamilies() []string

	// Return faces in a family and/or with a particular set of attributes
	GetFaces(family string, flags VGFontStyleFlags) []VGFace
}

// Abstract font face interface
type VGFace interface {
	// Get Face Name (from the filename)
	GetName() string

	// Get Face Index
	GetIndex() uint

	// Get Number of faces within the file
	GetNumFaces() uint

	// Number of glyphs for the face
	GetNumGlyphs() uint

	// Return name of font family
	GetFamily() string

	// Return style name of font face
	GetStyle() string
}

// Font face query flags
type VGFontStyleFlags uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Constants used for querying faces for VGFontDriver
	VG_FONT_STYLE_ANY     VGFontStyleFlags = iota
	VG_FONT_STYLE_REGULAR
	VG_FONT_STYLE_BOLD
	VG_FONT_STYLE_ITALIC
	VG_FONT_STYLE_BOLDITALIC
)

