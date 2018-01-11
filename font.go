/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"os"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Font face query flags
type FontStyleFlags uint16

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// Manager for fonts
type FontManager interface {
	Driver

	// Open a font face - first face at index 0 is loaded
	OpenFace(path string) (FontFace, error)

	// Open a font face - indexed within file of several faces
	OpenFaceAtIndex(path string, index uint) (FontFace, error)

	// Open font faces at path, checking to see if individual files should
	// be opened through a callback function
	OpenFacesAtPath(path string, callback func(path string, info os.FileInfo) bool) error

	// Destroy a font face
	DestroyFace(FontFace) error

	// Return an array of font families which are loaded
	GetFamilies() []string

	// Return faces in a family and/or with a particular set of attributes
	GetFaces(family string, flags FontStyleFlags) []FontFace
}

// Abstract font face interface
type FontFace interface {

	// Get Face Name (from the filename)
	Name() string

	// Get Face Index
	Index() uint

	// Get Number of faces within the file
	NumFaces() uint

	// Number of glyphs for the face
	NumGlyphs() uint

	// Return name of font family
	Family() string

	// Return style name of font face
	Style() string
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Constants used for querying faces for VGFontDriver
	FONT_STYLE_ANY FontStyleFlags = iota
	FONT_STYLE_REGULAR
	FONT_STYLE_BOLD
	FONT_STYLE_ITALIC
	FONT_STYLE_BOLDITALIC
)
