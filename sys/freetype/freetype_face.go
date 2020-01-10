// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype

import (
	"unsafe"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: freetype2
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// FACE FUNCTIONS

func FT_NewFace(handle FT_Library, path string, index uint) (FT_Face, error) {
	var face C.FT_Face
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	if err := FT_Error(C.FT_New_Face(handle, cstr, C.FT_Long(index), &face)); err != FT_SUCCESS {
		return nil, err
	} else {
		return FT_Face(face), nil
	}
}

func FT_DoneFace(handle FT_Face) error {
	if err := FT_Error(C.FT_Done_Face(handle)); err != FT_SUCCESS {
		return err
	} else {
		return nil
	}
}

func FT_SelectCharmap(handle FT_Face, encoding FT_Encoding) error {
	if err := FT_Error(C.FT_Select_Charmap(handle, C.FT_Encoding(encoding))); err != FT_SUCCESS {
		return err
	} else {
		return nil
	}
}

func FT_FaceFamily(handle FT_Face) string {
	return C.GoString(handle.family_name)
}

func FT_FaceStyle(handle FT_Face) string {
	return C.GoString(handle.style_name)
}

func FT_FaceIndex(handle FT_Face) uint {
	return uint(handle.face_index)
}

func FT_FaceNumFaces(handle FT_Face) uint {
	return uint(handle.num_faces)
}

func FT_FaceNumGlyphs(handle FT_Face) uint {
	return uint(handle.num_glyphs)
}

func FT_FaceStyleFlags(handle FT_Face) gopi.FontFlags {
	return gopi.FontFlags(handle.style_flags)
}

func FT_SetPixelSizes(handle FT_Face, size float32) error {
	if err := FT_Error(C.FT_Set_Pixel_Sizes(handle, 0, C.FT_UInt(size))); err != FT_SUCCESS {
		return err
	} else {
		return nil
	}
}

func FT_SetCharSize(handle FT_Face, points float32, ppi uint) error {
	if err := FT_Error(C.FT_Set_Char_Size(handle, 0, C.FT_F26Dot6(points*64.0), 0, C.FT_UInt(ppi))); err != FT_SUCCESS {
		return err
	} else {
		return nil
	}
}

// This method returns a bitmap for a rune. The returned values are a pointer
// to the bitmap pixels
func FT_Load_Glyph(handle FT_Face, value rune) (uintptr, error) {

	// Get Glyph
	glyph_index := C.FT_Get_Char_Index(handle, C.FT_ULong(value))
	if glyph_index == 0 {
		return 0, gopi.ErrBadParameter.WithPrefix("rune")
	}

	// Render Glyph
	if err := FT_Error(C.FT_Load_Glyph(handle, glyph_index, C.FT_Int32(FT_LOAD_RENDER))); err != FT_SUCCESS {
		return 0, err
	} else {
		return 0, nil
	}
	/*
		// Compute relevant information
		bitmap := this.handle.glyph.bitmap
		pixel_mode := VGFontBitmapPixelMode(bitmap.pixel_mode)
		size := khronos.EGLSize{Width: uint(bitmap.width), Height: uint(bitmap.rows)}
		advance := khronos.EGLSize{Width: uint(this.handle.glyph.advance.x >> 6), Height: uint(this.handle.glyph.advance.y >> 6)}
		stride := uint(bitmap.pitch)

		// Success
		return uintptr(unsafe.Pointer(bitmap.buffer)), pixel_mode, size, advance, stride, nil
	*/
}
