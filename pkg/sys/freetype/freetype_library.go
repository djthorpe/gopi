// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: freetype2
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// LIBRARY FUNCTIONS

func FT_Init() (FT_Library, error) {
	var handle C.FT_Library
	if err := FT_Error(C.FT_Init_FreeType((*C.FT_Library)(&handle))); err != FT_SUCCESS {
		return nil, err
	} else {
		return FT_Library(handle), nil
	}
}

func FT_Destroy(handle FT_Library) error {
	if err := FT_Error(C.FT_Done_FreeType(handle)); err != FT_SUCCESS {
		return err
	} else {
		return nil
	}
}

func FT_Library_Version(handle FT_Library) (int, int, int) {
	var major, minor, patch C.FT_Int
	C.FT_Library_Version(handle, (*C.FT_Int)(&major), (*C.FT_Int)(&minor), (*C.FT_Int)(&patch))
	return int(major), int(minor), int(patch)
}
