package gopi

import (
	"fmt"
	"os"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	FontFlags    uint16
	FontSizeUnit uint

	// SurfaceFlags are flags associated with renderable surface
	SurfaceFlags uint16

	// SurfaceFormat defines the pixel format for a surface
	SurfaceFormat uint
)

type FontSize struct {
	Size float32
	Unit FontSizeUnit
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// SurfaceManager to manage graphics surfaces
type SurfaceManager interface {
	CreateBackground(Display, SurfaceFlags) (Surface, error)
	DisposeSurface(Surface) error

	SwapBuffers() error
	/*
		CreateSurfaceWithBitmap(Bitmap, SurfaceFlags, float32, uint16, Point, Size) (Surface, error)
		CreateSurface(SurfaceFlags, float32, uint16, Point, Size) (Surface, error)
	*/
}

type Surface interface {
	Origin() Point
	Size() Size
}

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

	// Return faces in a family and/or with a particular set of attributes
	Faces(family string, flags FontFlags) []FontFace
}

// FontFace represents a typeface
type FontFace interface {
	Name() string     // Get Face Name (from the filename)
	Index() uint      // Get Face Index
	NumFaces() uint   // Get Number of faces within the file
	NumGlyphs() uint  // Number of glyphs for the face
	Family() string   // Return name of font family
	Style() string    // Return style name of font face
	Flags() FontFlags // Return properties for face
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

const (
	FONT_SIZE_PIXELS FontSizeUnit = iota
	FONT_SIZE_POINTS
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SURFACE_FLAG_BITMAP SurfaceFlags = (1 << iota)
	SURFACE_FLAG_OPENGL
	SURFACE_FLAG_OPENGL_ES
	SURFACE_FLAG_OPENGL_ES2
	SURFACE_FLAG_OPENVG
	SURFACE_FLAG_NONE SurfaceFlags = 0
	SURFACE_FLAG_MIN               = SURFACE_FLAG_BITMAP
	SURFACE_FLAG_MAX               = SURFACE_FLAG_OPENVG
)

const (
	SURFACE_FMT_NONE   SurfaceFormat = iota
	SURFACE_FMT_RGBA32               // 4 bytes per pixel with transparency
	SURFACE_FMT_XRGB32               // 4 bytes per pixel without transparency
	SURFACE_FMT_RGB888               // 3 bytes per pixel
	SURFACE_FMT_RGB565               // 2 bytes per pixel
	SURFACE_FMT_1BPP                 // 1 bit per pixel (Mono)
	SURFACE_FMT_MAX    = SURFACE_FMT_1BPP
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

func (f SurfaceFlags) String() string {
	if f == SURFACE_FLAG_NONE {
		return f.StringFlag()
	}
	str := ""
	for v := SURFACE_FLAG_MIN; v <= SURFACE_FLAG_MAX; v <<= 1 {
		if f&v == v {
			str += v.StringFlag() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (f SurfaceFlags) StringFlag() string {
	switch f {
	case SURFACE_FLAG_NONE:
		return "SURFACE_FLAG_NONE"
	case SURFACE_FLAG_BITMAP:
		return "SURFACE_FLAG_BITMAP"
	case SURFACE_FLAG_OPENGL:
		return "SURFACE_FLAG_OPENGL"
	case SURFACE_FLAG_OPENGL_ES:
		return "SURFACE_FLAG_OPENGL_ES"
	case SURFACE_FLAG_OPENGL_ES2:
		return "SURFACE_FLAG_OPENGL_ES2"
	case SURFACE_FLAG_OPENVG:
		return "SURFACE_FLAG_OPENVG"
	default:
		return "[?? Invalid SurfaceFlags value]"
	}
}

func (f SurfaceFormat) String() string {
	switch f {
	case SURFACE_FMT_NONE:
		return "SURFACE_FMT_NONE"
	case SURFACE_FMT_RGBA32:
		return "SURFACE_FMT_RGBA32"
	case SURFACE_FMT_XRGB32:
		return "SURFACE_FMT_XRGB32"
	case SURFACE_FMT_RGB888:
		return "SURFACE_FMT_RGB888"
	case SURFACE_FMT_RGB565:
		return "SURFACE_FMT_RGB565"
	default:
		return "[?? Invalid SurfaceFormat value]"
	}
}
