// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"fmt"
	"os"
	"sync"
	"unsafe"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
  #cgo CFLAGS:   -I/usr/include/freetype2
  #cgo LDFLAGS:  -lfreetype
  #include <ft2build.h>
  #include FT_FREETYPE_H
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type FontManager struct {
	// RootPath for relative OpenFace calls
	RootPath string
}

type manager struct {
	log                 gopi.Logger
	root_path           string
	lock                sync.Mutex
	library             ftLibrary
	major, minor, patch ftInt
}

type face struct {
}

type (
	ftError   C.FT_Error
	ftLibrary uintptr
	ftInt     int
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FT_SUCCESS ftError = 0

	/* generic errors */
	FT_ERROR_Cannot_Open_Resource  = 0x01
	FT_ERROR_Unknown_File_Format   = 0x02
	FT_ERROR_Invalid_File_Format   = 0x03
	FT_ERROR_Invalid_Version       = 0x04
	FT_ERROR_Lower_Module_Version  = 0x05
	FT_ERROR_Invalid_Argument      = 0x06
	FT_ERROR_Unimplemented_Feature = 0x07
	FT_ERROR_Invalid_Table         = 0x08
	FT_ERROR_Invalid_Offset        = 0x09
	FT_ERROR_Array_Too_Large       = 0x0A
	FT_ERROR_Missing_Module        = 0x0B
	FT_ERROR_Missing_Property      = 0x0C

	/* glyph/character errors */
	FT_ERROR_Invalid_Glyph_Index    = 0x10
	FT_ERROR_Invalid_Character_Code = 0x11
	FT_ERROR_Invalid_Glyph_Format   = 0x12
	FT_ERROR_Cannot_Render_Glyph    = 0x13
	FT_ERROR_Invalid_Outline        = 0x14
	FT_ERROR_Invalid_Composite      = 0x15
	FT_ERROR_Too_Many_Hints         = 0x16
	FT_ERROR_Invalid_Pixel_Size     = 0x17

	/* handle errors */
	FT_ERROR_Invalid_Handle         = 0x20
	FT_ERROR_Invalid_Library_Handle = 0x21
	FT_ERROR_Invalid_Driver_Handle  = 0x22
	FT_ERROR_Invalid_Face_Handle    = 0x23
	FT_ERROR_Invalid_Size_Handle    = 0x24
	FT_ERROR_Invalid_Slot_Handle    = 0x25
	FT_ERROR_Invalid_CharMap_Handle = 0x26
	FT_ERROR_Invalid_Cache_Handle   = 0x27
	FT_ERROR_Invalid_Stream_Handle  = 0x28

	/* driver errors */
	FT_ERROR_Too_Many_Drivers    = 0x30
	FT_ERROR_Too_Many_Extensions = 0x31

	/* memory errors */
	FT_ERROR_Out_Of_Memory   = 0x40
	FT_ERROR_Unlisted_Object = 0x41

	/* stream errors */
	FT_ERROR_Cannot_Open_Stream       = 0x51
	FT_ERROR_Invalid_Stream_Seek      = 0x52
	FT_ERROR_Invalid_Stream_Skip      = 0x53
	FT_ERROR_Invalid_Stream_Read      = 0x54
	FT_ERROR_Invalid_Stream_Operation = 0x55
	FT_ERROR_Invalid_Frame_Operation  = 0x56
	FT_ERROR_Nested_Frame_Access      = 0x57
	FT_ERROR_Invalid_Frame_Read       = 0x58

	/* raster errors */
	FT_ERROR_Raster_Uninitialized   = 0x60
	FT_ERROR_Raster_Corrupted       = 0x61
	FT_ERROR_Raster_Overflow        = 0x62
	FT_ERROR_Raster_Negative_Height = 0x63

	/* cache errors */
	FT_ERROR_Too_Many_Caches = 0x70

	/* TrueType and SFNT errors */
	FT_ERROR_Invalid_Opcode            = 0x80
	FT_ERROR_Too_Few_Arguments         = 0x81
	FT_ERROR_Stack_Overflow            = 0x82
	FT_ERROR_Code_Overflow             = 0x83
	FT_ERROR_Bad_Argument              = 0x84
	FT_ERROR_Divide_By_Zero            = 0x85
	FT_ERROR_Invalid_Reference         = 0x86
	FT_ERROR_Debug_OpCode              = 0x87
	FT_ERROR_ENDF_In_Exec_Stream       = 0x88
	FT_ERROR_Nested_DEFS               = 0x89
	FT_ERROR_Invalid_CodeRange         = 0x8A
	FT_ERROR_Execution_Too_Long        = 0x8B
	FT_ERROR_Too_Many_Function_Defs    = 0x8C
	FT_ERROR_Too_Many_Instruction_Defs = 0x8D
	FT_ERROR_Table_Missing             = 0x8E
	FT_ERROR_Horiz_Header_Missing      = 0x8F
	FT_ERROR_Locations_Missing         = 0x90
	FT_ERROR_Name_Table_Missing        = 0x91
	FT_ERROR_CMap_Table_Missing        = 0x92
	FT_ERROR_Hmtx_Table_Missing        = 0x93
	FT_ERROR_Post_Table_Missing        = 0x94
	FT_ERROR_Invalid_Horiz_Metrics     = 0x95
	FT_ERROR_Invalid_CharMap_Format    = 0x96
	FT_ERROR_Invalid_PPem              = 0x97
	FT_ERROR_Invalid_Vert_Metrics      = 0x98
	FT_ERROR_Could_Not_Find_Context    = 0x99
	FT_ERROR_Invalid_Post_Table_Format = 0x9A
	FT_ERROR_Invalid_Post_Table        = 0x9B

	/* CFF CID and Type 1 errors */
	FT_ERROR_Syntax_Error          = 0xA0
	FT_ERROR_Stack_Underflow       = 0xA1
	FT_ERROR_Ignore                = 0xA2
	FT_ERROR_No_Unicode_Glyph_Name = 0xA3
	FT_ERROR_Glyph_Too_Big         = 0xA4

	/* BDF errors */
	FT_ERROR_Missing_Startfont_Field       = 0xB0
	FT_ERROR_Missing_Font_Field            = 0xB1
	FT_ERROR_Missing_Size_Field            = 0xB2
	FT_ERROR_Missing_Fontboundingbox_Field = 0xB3
	FT_ERROR_Missing_Chars_Field           = 0xB4
	FT_ERROR_Missing_Startchar_Field       = 0xB5
	FT_ERROR_Missing_Encoding_Field        = 0xB6
	FT_ERROR_Missing_Bbx_Field             = 0xB7
	FT_ERROR_Bbx_Too_Big                   = 0xB8
	FT_ERROR_Corrupted_Font_Header         = 0xB9
	FT_ERROR_Corrupted_Font_Glyphs         = 0xBA
)

const (
	FT_NO_LIBRARY = ftLibrary(0)
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config FontManager) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.font.rpi.FontManager.Open>{ root_path=%v }", config.RootPath)
	if config.RootPath != "" {
		if stat, err := os.Stat(config.RootPath); os.IsNotExist(err) || stat.IsDir() == false {
			return nil, gopi.ErrBadParameter
		} else if err != nil {
			return nil, err
		}
	}
	this := new(manager)
	this.log = log
	this.root_path = config.RootPath

	this.lock.Lock()
	defer this.lock.Unlock()

	if library, err := ftInit(); err != FT_SUCCESS {
		return nil, os.NewSyscallError("ftInit", err)
	} else {
		this.library = library
	}

	this.major, this.minor, this.patch = ftLibraryVersion(this.library)

	return this, nil
}

func (this *manager) Close() error {
	this.log.Debug("<sys.font.rpi.FontManager.Close>{ handle=0x%X }", this.library)

	this.lock.Lock()
	defer this.lock.Unlock()

	if this.library == FT_NO_LIBRARY {
		return nil
	} else if err := ftDestroy(this.library); err != FT_SUCCESS {
		this.library = FT_NO_LIBRARY
		return err
	} else {
		this.library = FT_NO_LIBRARY
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *manager) String() string {
	return fmt.Sprintf("<sys.font.rpi.FontManager>{ handle=0x%X version={%v,%v,%v} }", this.library, this.major, this.minor, this.patch)
}

func (e ftError) Error() string {
	switch e {
	case FT_SUCCESS:
		return "FT_SUCCESS"
	case FT_ERROR_Cannot_Open_Resource:
		return "Cannot_Open_Resource "
	case FT_ERROR_Unknown_File_Format:
		return "Unknown_File_Format "
	case FT_ERROR_Invalid_File_Format:
		return "Invalid_File_Format "
	case FT_ERROR_Invalid_Version:
		return "Invalid_Version "
	case FT_ERROR_Lower_Module_Version:
		return "Lower_Module_Version "
	case FT_ERROR_Invalid_Argument:
		return "Invalid_Argument "
	case FT_ERROR_Unimplemented_Feature:
		return "Unimplemented_Feature "
	case FT_ERROR_Invalid_Table:
		return "Invalid_Table "
	case FT_ERROR_Invalid_Offset:
		return "Invalid_Offset "
	case FT_ERROR_Array_Too_Large:
		return "Array_Too_Large "
	case FT_ERROR_Missing_Module:
		return "Missing_Module "
	case FT_ERROR_Missing_Property:
		return "Missing_Property "
	case FT_ERROR_Invalid_Glyph_Index:
		return "Invalid_Glyph_Index "
	case FT_ERROR_Invalid_Character_Code:
		return "Invalid_Character_Code "
	case FT_ERROR_Invalid_Glyph_Format:
		return "Invalid_Glyph_Format "
	case FT_ERROR_Cannot_Render_Glyph:
		return "Cannot_Render_Glyph "
	case FT_ERROR_Invalid_Outline:
		return "Invalid_Outline "
	case FT_ERROR_Invalid_Composite:
		return "Invalid_Composite "
	case FT_ERROR_Too_Many_Hints:
		return "Too_Many_Hints "
	case FT_ERROR_Invalid_Pixel_Size:
		return "Invalid_Pixel_Size "
	case FT_ERROR_Invalid_Handle:
		return "Invalid_Handle "
	case FT_ERROR_Invalid_Library_Handle:
		return "Invalid_Library_Handle "
	case FT_ERROR_Invalid_Driver_Handle:
		return "Invalid_Driver_Handle "
	case FT_ERROR_Invalid_Face_Handle:
		return "Invalid_Face_Handle "
	case FT_ERROR_Invalid_Size_Handle:
		return "Invalid_Size_Handle "
	case FT_ERROR_Invalid_Slot_Handle:
		return "Invalid_Slot_Handle "
	case FT_ERROR_Invalid_CharMap_Handle:
		return "Invalid_CharMap_Handle "
	case FT_ERROR_Invalid_Cache_Handle:
		return "Invalid_Cache_Handle "
	case FT_ERROR_Invalid_Stream_Handle:
		return "Invalid_Stream_Handle "
	case FT_ERROR_Too_Many_Drivers:
		return "Too_Many_Drivers "
	case FT_ERROR_Too_Many_Extensions:
		return "Too_Many_Extensions "
	case FT_ERROR_Out_Of_Memory:
		return "Out_Of_Memory "
	case FT_ERROR_Unlisted_Object:
		return "Unlisted_Object "
	case FT_ERROR_Cannot_Open_Stream:
		return "Cannot_Open_Stream "
	case FT_ERROR_Invalid_Stream_Seek:
		return "Invalid_Stream_Seek "
	case FT_ERROR_Invalid_Stream_Skip:
		return "Invalid_Stream_Skip "
	case FT_ERROR_Invalid_Stream_Read:
		return "Invalid_Stream_Read "
	case FT_ERROR_Invalid_Stream_Operation:
		return "Invalid_Stream_Operation "
	case FT_ERROR_Invalid_Frame_Operation:
		return "Invalid_Frame_Operation "
	case FT_ERROR_Nested_Frame_Access:
		return "Nested_Frame_Access "
	case FT_ERROR_Invalid_Frame_Read:
		return "Invalid_Frame_Read "
	case FT_ERROR_Raster_Uninitialized:
		return "Raster_Uninitialized "
	case FT_ERROR_Raster_Corrupted:
		return "Raster_Corrupted "
	case FT_ERROR_Raster_Overflow:
		return "Raster_Overflow "
	case FT_ERROR_Raster_Negative_Height:
		return "Raster_Negative_Height "
	case FT_ERROR_Too_Many_Caches:
		return "Too_Many_Caches "
	case FT_ERROR_Invalid_Opcode:
		return "Invalid_Opcode "
	case FT_ERROR_Too_Few_Arguments:
		return "Too_Few_Arguments "
	case FT_ERROR_Stack_Overflow:
		return "Stack_Overflow "
	case FT_ERROR_Code_Overflow:
		return "Code_Overflow "
	case FT_ERROR_Bad_Argument:
		return "Bad_Argument "
	case FT_ERROR_Divide_By_Zero:
		return "Divide_By_Zero "
	case FT_ERROR_Invalid_Reference:
		return "Invalid_Reference "
	case FT_ERROR_Debug_OpCode:
		return "Debug_OpCode "
	case FT_ERROR_ENDF_In_Exec_Stream:
		return "ENDF_In_Exec_Stream "
	case FT_ERROR_Nested_DEFS:
		return "Nested_DEFS "
	case FT_ERROR_Invalid_CodeRange:
		return "Invalid_CodeRange "
	case FT_ERROR_Execution_Too_Long:
		return "Execution_Too_Long "
	case FT_ERROR_Too_Many_Function_Defs:
		return "Too_Many_Function_Defs "
	case FT_ERROR_Too_Many_Instruction_Defs:
		return "Too_Many_Instruction_Defs "
	case FT_ERROR_Table_Missing:
		return "Table_Missing "
	case FT_ERROR_Horiz_Header_Missing:
		return "Horiz_Header_Missing "
	case FT_ERROR_Locations_Missing:
		return "Locations_Missing "
	case FT_ERROR_Name_Table_Missing:
		return "Name_Table_Missing "
	case FT_ERROR_CMap_Table_Missing:
		return "CMap_Table_Missing "
	case FT_ERROR_Hmtx_Table_Missing:
		return "Hmtx_Table_Missing "
	case FT_ERROR_Post_Table_Missing:
		return "Post_Table_Missing "
	case FT_ERROR_Invalid_Horiz_Metrics:
		return "Invalid_Horiz_Metrics "
	case FT_ERROR_Invalid_CharMap_Format:
		return "Invalid_CharMap_Format "
	case FT_ERROR_Invalid_PPem:
		return "Invalid_PPem "
	case FT_ERROR_Invalid_Vert_Metrics:
		return "Invalid_Vert_Metrics "
	case FT_ERROR_Could_Not_Find_Context:
		return "Could_Not_Find_Context "
	case FT_ERROR_Invalid_Post_Table_Format:
		return "Invalid_Post_Table_Format "
	case FT_ERROR_Invalid_Post_Table:
		return "Invalid_Post_Table "
	case FT_ERROR_Syntax_Error:
		return "Syntax_Error "
	case FT_ERROR_Stack_Underflow:
		return "Stack_Underflow "
	case FT_ERROR_Ignore:
		return "Ignore "
	case FT_ERROR_No_Unicode_Glyph_Name:
		return "No_Unicode_Glyph_Name "
	case FT_ERROR_Glyph_Too_Big:
		return "Glyph_Too_Big "
	case FT_ERROR_Missing_Startfont_Field:
		return "Missing_Startfont_Field "
	case FT_ERROR_Missing_Font_Field:
		return "Missing_Font_Field "
	case FT_ERROR_Missing_Size_Field:
		return "Missing_Size_Field "
	case FT_ERROR_Missing_Fontboundingbox_Field:
		return "Missing_Fontboundingbox_Field "
	case FT_ERROR_Missing_Chars_Field:
		return "Missing_Chars_Field "
	case FT_ERROR_Missing_Startchar_Field:
		return "Missing_Startchar_Field "
	case FT_ERROR_Missing_Encoding_Field:
		return "Missing_Encoding_Field "
	case FT_ERROR_Missing_Bbx_Field:
		return "Missing_Bbx_Field "
	case FT_ERROR_Bbx_Too_Big:
		return "Bbx_Too_Big "
	case FT_ERROR_Corrupted_Font_Header:
		return "Corrupted_Font_Header "
	case FT_ERROR_Corrupted_Font_Glyphs:
		return "Corrupted_Font_Glyphsft"
	default: // TODO Widen out error messages!
		return "FT_ERROR - TODO: Add error codes"
	}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

func (this *manager) OpenFace(path string) (gopi.FontFace, error) {
	return this.OpenFaceAtIndex(path, 0)
}

func (this *manager) OpenFaceAtIndex(path string, index uint) (gopi.FontFace, error) {
	this.log.Debug2("<sys.font.rpi.FontManager.OpenFaceAtIndex{ path=%v index=%v }", path, index)
	return nil, gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS

func ftInit() (ftLibrary, ftError) {
	var handle C.FT_Library
	if err := C.FT_Init_FreeType((*C.FT_Library)(&handle)); ftError(err) != FT_SUCCESS {
		return FT_NO_LIBRARY, ftError(err)
	} else {
		return ftLibrary(unsafe.Pointer(handle)), FT_SUCCESS
	}
}

func ftDestroy(handle ftLibrary) ftError {
	if err := C.FT_Done_FreeType((C.FT_Library)(unsafe.Pointer(handle))); ftError(err) != FT_SUCCESS {
		return ftError(err)
	} else {
		return FT_SUCCESS
	}
}

func ftLibraryVersion(handle ftLibrary) (ftInt, ftInt, ftInt) {
	var major, minor, patch ftInt
	C.FT_Library_Version((C.FT_Library)(unsafe.Pointer(handle)), (*C.FT_Int)(unsafe.Pointer(&major)), (*C.FT_Int)(unsafe.Pointer(&minor)), (*C.FT_Int)(unsafe.Pointer(&patch)))
	return major, minor, patch
}
