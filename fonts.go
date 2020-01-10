/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"os"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Font flags
type FontFlags uint16

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// FontManager for font management
type FontManager interface {

	// Open a font face - first face at index 0 is loaded
	OpenFace(path string) (FontFace, error)

	// Open a font face - indexed within file of several faces
	OpenFaceAtIndex(path string, index uint) (FontFace, error)

	// Open font faces at path, checking to see if individual files should
	// be opened through a callback function
	OpenFacesAtPath(path string, callback func(manager FontManager, path string, info os.FileInfo) bool) error

	// Destroy a font face
	DestroyFace(FontFace) error

	// Return an array of font families which are loaded
	Families() []string

	// Return open face for filepath
	FaceForPath(path string) FontFace

	// Return faces in a family and/or with a particular set of attributes
	Faces(family string, flags FontFlags) []FontFace

	// Implements gopi.Unit
	Unit
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

	// Return properties for face
	Flags() FontFlags
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FONT_FLAGS_NONE             FontFlags = 0x0000
	FONT_FLAGS_STYLE_ITALIC     FontFlags = 0x0001
	FONT_FLAGS_STYLE_BOLD       FontFlags = 0x0002
	FONT_FLAGS_STYLE_BOLDITALIC FontFlags = 0x0003
	FONT_FLAGS_STYLE_REGULAR    FontFlags = 0x0004
	FONT_FLAGS_STYLE_ANY        FontFlags = 0x0007
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f FontFlags) String() string {
	if f == FONT_FLAGS_NONE {
		return f.StringFlag()
	}
	str := ""
	for v := FONT_FLAGS_STYLE_ITALIC; v <= FONT_FLAGS_STYLE_REGULAR; v <<= 1 {
		if f&v == v {
			str += v.StringFlag() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (f FontFlags) StringFlag() string {
	switch f {
	case FONT_FLAGS_NONE:
		return "FONT_FLAGS_NONE"
	case FONT_FLAGS_STYLE_REGULAR:
		return "FONT_FLAGS_STYLE_REGULAR"
	case FONT_FLAGS_STYLE_BOLD:
		return "FONT_FLAGS_STYLE_BOLD"
	case FONT_FLAGS_STYLE_ITALIC:
		return "FONT_FLAGS_STYLE_ITALIC"
	case FONT_FLAGS_STYLE_BOLDITALIC:
		return "FONT_FLAGS_STYLE_BOLDITALIC"
	case FONT_FLAGS_STYLE_ANY:
		return "FONT_FLAGS_STYLE_ANY"
	default:
		return fmt.Sprintf("[?? Invalid FontFlags value %d]", int(f))
	}
}
