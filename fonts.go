/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import "os"

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Font flags
type FontFlags uint16

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// FontManager for font management
type FontManager interface {
	Driver

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
	FONT_FLAGS_STYLE_REGULAR    FontFlags = 0x0001
	FONT_FLAGS_STYLE_BOLD       FontFlags = 0x0002
	FONT_FLAGS_STYLE_ITALIC     FontFlags = 0x0004
	FONT_FLAGS_STYLE_BOLDITALIC FontFlags = 0x0006
	FONT_FLAGS_STYLE_ANY        FontFlags = 0x0007
)
