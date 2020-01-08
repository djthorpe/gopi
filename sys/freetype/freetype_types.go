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
// TYPES

type (
	FT_Error     C.FT_Error
	FT_Library   C.FT_Library
	FT_Face      C.FT_Face
	FT_Encoding  C.FT_Encoding
	FT_LoadFlags uint32
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FT_SUCCESS FT_Error = 0

	FT_ERROR_Cannot_Open_Resource          FT_Error = 0x01 // Generic errors
	FT_ERROR_Unknown_File_Format           FT_Error = 0x02
	FT_ERROR_Invalid_File_Format           FT_Error = 0x03
	FT_ERROR_Invalid_Version               FT_Error = 0x04
	FT_ERROR_Lower_Module_Version          FT_Error = 0x05
	FT_ERROR_Invalid_Argument              FT_Error = 0x06
	FT_ERROR_Unimplemented_Feature         FT_Error = 0x07
	FT_ERROR_Invalid_Table                 FT_Error = 0x08
	FT_ERROR_Invalid_Offset                FT_Error = 0x09
	FT_ERROR_Array_Too_Large               FT_Error = 0x0A
	FT_ERROR_Missing_Module                FT_Error = 0x0B
	FT_ERROR_Missing_Property              FT_Error = 0x0C
	FT_ERROR_Invalid_Glyph_Index           FT_Error = 0x10 // Glyph/character errors
	FT_ERROR_Invalid_Character_Code        FT_Error = 0x11
	FT_ERROR_Invalid_Glyph_Format          FT_Error = 0x12
	FT_ERROR_Cannot_Render_Glyph           FT_Error = 0x13
	FT_ERROR_Invalid_Outline               FT_Error = 0x14
	FT_ERROR_Invalid_Composite             FT_Error = 0x15
	FT_ERROR_Too_Many_Hints                FT_Error = 0x16
	FT_ERROR_Invalid_Pixel_Size            FT_Error = 0x17
	FT_ERROR_Invalid_Handle                FT_Error = 0x20 // Handle errors
	FT_ERROR_Invalid_Library_Handle        FT_Error = 0x21
	FT_ERROR_Invalid_Driver_Handle         FT_Error = 0x22
	FT_ERROR_Invalid_Face_Handle           FT_Error = 0x23
	FT_ERROR_Invalid_Size_Handle           FT_Error = 0x24
	FT_ERROR_Invalid_Slot_Handle           FT_Error = 0x25
	FT_ERROR_Invalid_CharMap_Handle        FT_Error = 0x26
	FT_ERROR_Invalid_Cache_Handle          FT_Error = 0x27
	FT_ERROR_Invalid_Stream_Handle         FT_Error = 0x28
	FT_ERROR_Too_Many_Drivers              FT_Error = 0x30 /* driver errors */
	FT_ERROR_Too_Many_Extensions           FT_Error = 0x31
	FT_ERROR_Out_Of_Memory                 FT_Error = 0x40 /* memory errors */
	FT_ERROR_Unlisted_Object               FT_Error = 0x41
	FT_ERROR_Cannot_Open_Stream            FT_Error = 0x51 /* stream errors */
	FT_ERROR_Invalid_Stream_Seek           FT_Error = 0x52
	FT_ERROR_Invalid_Stream_Skip           FT_Error = 0x53
	FT_ERROR_Invalid_Stream_Read           FT_Error = 0x54
	FT_ERROR_Invalid_Stream_Operation      FT_Error = 0x55
	FT_ERROR_Invalid_Frame_Operation       FT_Error = 0x56
	FT_ERROR_Nested_Frame_Access           FT_Error = 0x57
	FT_ERROR_Invalid_Frame_Read            FT_Error = 0x58
	FT_ERROR_Raster_Uninitialized          FT_Error = 0x60 /* raster errors */
	FT_ERROR_Raster_Corrupted              FT_Error = 0x61
	FT_ERROR_Raster_Overflow               FT_Error = 0x62
	FT_ERROR_Raster_Negative_Height        FT_Error = 0x63
	FT_ERROR_Too_Many_Caches               FT_Error = 0x70 /* cache errors */
	FT_ERROR_Invalid_Opcode                FT_Error = 0x80 /* TrueType and SFNT errors */
	FT_ERROR_Too_Few_Arguments             FT_Error = 0x81
	FT_ERROR_Stack_Overflow                FT_Error = 0x82
	FT_ERROR_Code_Overflow                 FT_Error = 0x83
	FT_ERROR_Bad_Argument                  FT_Error = 0x84
	FT_ERROR_Divide_By_Zero                FT_Error = 0x85
	FT_ERROR_Invalid_Reference             FT_Error = 0x86
	FT_ERROR_Debug_OpCode                  FT_Error = 0x87
	FT_ERROR_ENDF_In_Exec_Stream           FT_Error = 0x88
	FT_ERROR_Nested_DEFS                   FT_Error = 0x89
	FT_ERROR_Invalid_CodeRange             FT_Error = 0x8A
	FT_ERROR_Execution_Too_Long            FT_Error = 0x8B
	FT_ERROR_Too_Many_Function_Defs        FT_Error = 0x8C
	FT_ERROR_Too_Many_Instruction_Defs     FT_Error = 0x8D
	FT_ERROR_Table_Missing                 FT_Error = 0x8E
	FT_ERROR_Horiz_Header_Missing          FT_Error = 0x8F
	FT_ERROR_Locations_Missing             FT_Error = 0x90
	FT_ERROR_Name_Table_Missing            FT_Error = 0x91
	FT_ERROR_CMap_Table_Missing            FT_Error = 0x92
	FT_ERROR_Hmtx_Table_Missing            FT_Error = 0x93
	FT_ERROR_Post_Table_Missing            FT_Error = 0x94
	FT_ERROR_Invalid_Horiz_Metrics         FT_Error = 0x95
	FT_ERROR_Invalid_CharMap_Format        FT_Error = 0x96
	FT_ERROR_Invalid_PPem                  FT_Error = 0x97
	FT_ERROR_Invalid_Vert_Metrics          FT_Error = 0x98
	FT_ERROR_Could_Not_Find_Context        FT_Error = 0x99
	FT_ERROR_Invalid_Post_Table_Format     FT_Error = 0x9A
	FT_ERROR_Invalid_Post_Table            FT_Error = 0x9B
	FT_ERROR_Syntax_Error                  FT_Error = 0xA0 /* CFF CID and Type 1 errors */
	FT_ERROR_Stack_Underflow               FT_Error = 0xA1
	FT_ERROR_Ignore                        FT_Error = 0xA2
	FT_ERROR_No_Unicode_Glyph_Name         FT_Error = 0xA3
	FT_ERROR_Glyph_Too_Big                 FT_Error = 0xA4
	FT_ERROR_Missing_Startfont_Field       FT_Error = 0xB0 /* BDF errors */
	FT_ERROR_Missing_Font_Field            FT_Error = 0xB1
	FT_ERROR_Missing_Size_Field            FT_Error = 0xB2
	FT_ERROR_Missing_Fontboundingbox_Field FT_Error = 0xB3
	FT_ERROR_Missing_Chars_Field           FT_Error = 0xB4
	FT_ERROR_Missing_Startchar_Field       FT_Error = 0xB5
	FT_ERROR_Missing_Encoding_Field        FT_Error = 0xB6
	FT_ERROR_Missing_Bbx_Field             FT_Error = 0xB7
	FT_ERROR_Bbx_Too_Big                   FT_Error = 0xB8
	FT_ERROR_Corrupted_Font_Header         FT_Error = 0xB9
	FT_ERROR_Corrupted_Font_Glyphs         FT_Error = 0xBA
	FT_ERROR_MIN                           FT_Error = FT_ERROR_Cannot_Open_Resource // Min and max error values
	FT_ERROR_MAX                           FT_Error = FT_ERROR_Corrupted_Font_Glyphs
)

const (
	FT_LOAD_DEFAULT                     FT_LoadFlags = 0
	FT_LOAD_NO_SCALE                    FT_LoadFlags = (1 << 0)
	FT_LOAD_NO_HINTING                  FT_LoadFlags = (1 << 1)
	FT_LOAD_RENDER                      FT_LoadFlags = (1 << 2)
	FT_LOAD_NO_BITMAP                   FT_LoadFlags = (1 << 3)
	FT_LOAD_VERTICAL_LAYOUT             FT_LoadFlags = (1 << 4)
	FT_LOAD_FORCE_AUTOHINT              FT_LoadFlags = (1 << 5)
	FT_LOAD_CROP_BITMAP                 FT_LoadFlags = (1 << 6)
	FT_LOAD_PEDANTIC                    FT_LoadFlags = (1 << 7)
	FT_LOAD_IGNORE_GLOBAL_ADVANCE_WIDTH FT_LoadFlags = (1 << 9)
	FT_LOAD_NO_RECURSE                  FT_LoadFlags = (1 << 10)
	FT_LOAD_IGNORE_TRANSFORM            FT_LoadFlags = (1 << 11)
	FT_LOAD_MONOCHROME                  FT_LoadFlags = (1 << 12)
	FT_LOAD_LINEAR_DESIGN               FT_LoadFlags = (1 << 13)
	FT_LOAD_NO_AUTOHINT                 FT_LoadFlags = (1 << 15)
	/* Bits 16-19 are used by `FT_LOAD_TARGET_` */
	FT_LOAD_COLOR               FT_LoadFlags = (1 << 20)
	FT_LOAD_COMPUTE_METRICS     FT_LoadFlags = (1 << 21)
	FT_LOAD_BITMAP_METRICS_ONLY FT_LoadFlags = (1 << 22)
)

var (
	FT_ENCODING_UNICODE = FOURCC('u', 'n', 'i', 'c')
)

func FOURCC(a, b, c, d uint8) FT_Encoding {
	return FT_Encoding(uint32(d) | uint32(c)<<8 | uint32(b)<<16 | uint32(a)<<24)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e FT_Error) Error() string {
	switch e {
	case FT_SUCCESS:
		return "FT_SUCCESS"
	case FT_ERROR_Cannot_Open_Resource:
		return "FT_ERROR_Cannot_Open_Resource"
	case FT_ERROR_Unknown_File_Format:
		return "FT_ERROR_Unknown_File_Format"
	case FT_ERROR_Invalid_File_Format:
		return "FT_ERROR_Invalid_File_Format"
	case FT_ERROR_Invalid_Version:
		return "FT_ERROR_Invalid_Version"
	case FT_ERROR_Lower_Module_Version:
		return "FT_ERROR_Lower_Module_Version"
	case FT_ERROR_Invalid_Argument:
		return "FT_ERROR_Invalid_Argument"
	case FT_ERROR_Unimplemented_Feature:
		return "FT_ERROR_Unimplemented_Feature"
	case FT_ERROR_Invalid_Table:
		return "FT_ERROR_Invalid_Table"
	case FT_ERROR_Invalid_Offset:
		return "FT_ERROR_Invalid_Offset"
	case FT_ERROR_Array_Too_Large:
		return "FT_ERROR_Array_Too_Large"
	case FT_ERROR_Missing_Module:
		return "FT_ERROR_Missing_Module"
	case FT_ERROR_Missing_Property:
		return "FT_ERROR_Missing_Property"
	case FT_ERROR_Invalid_Glyph_Index:
		return "FT_ERROR_Invalid_Glyph_Index"
	case FT_ERROR_Invalid_Character_Code:
		return "FT_ERROR_Invalid_Character_Code"
	case FT_ERROR_Invalid_Glyph_Format:
		return "FT_ERROR_Invalid_Glyph_Format"
	case FT_ERROR_Cannot_Render_Glyph:
		return "FT_ERROR_Cannot_Render_Glyph"
	case FT_ERROR_Invalid_Outline:
		return "FT_ERROR_Invalid_Outline"
	case FT_ERROR_Invalid_Composite:
		return "FT_ERROR_Invalid_Composite"
	case FT_ERROR_Too_Many_Hints:
		return "FT_ERROR_Too_Many_Hints"
	case FT_ERROR_Invalid_Pixel_Size:
		return "FT_ERROR_Invalid_Pixel_Size"
	case FT_ERROR_Invalid_Handle:
		return "FT_ERROR_Invalid_Handle"
	case FT_ERROR_Invalid_Library_Handle:
		return "FT_ERROR_Invalid_Library_Handle"
	case FT_ERROR_Invalid_Driver_Handle:
		return "FT_ERROR_Invalid_Driver_Handle"
	case FT_ERROR_Invalid_Face_Handle:
		return "FT_ERROR_Invalid_Face_Handle"
	case FT_ERROR_Invalid_Size_Handle:
		return "FT_ERROR_Invalid_Size_Handle"
	case FT_ERROR_Invalid_Slot_Handle:
		return "FT_ERROR_Invalid_Slot_Handle"
	case FT_ERROR_Invalid_CharMap_Handle:
		return "FT_ERROR_Invalid_CharMap_Handle"
	case FT_ERROR_Invalid_Cache_Handle:
		return "FT_ERROR_Invalid_Cache_Handle"
	case FT_ERROR_Invalid_Stream_Handle:
		return "FT_ERROR_Invalid_Stream_Handle"
	case FT_ERROR_Too_Many_Drivers:
		return "FT_ERROR_Too_Many_Drivers"
	case FT_ERROR_Too_Many_Extensions:
		return "FT_ERROR_Too_Many_Extensions"
	case FT_ERROR_Out_Of_Memory:
		return "FT_ERROR_Out_Of_Memory"
	case FT_ERROR_Unlisted_Object:
		return "FT_ERROR_Unlisted_Object"
	case FT_ERROR_Cannot_Open_Stream:
		return "FT_ERROR_Cannot_Open_Stream"
	case FT_ERROR_Invalid_Stream_Seek:
		return "FT_ERROR_Invalid_Stream_Seek"
	case FT_ERROR_Invalid_Stream_Skip:
		return "FT_ERROR_Invalid_Stream_Skip"
	case FT_ERROR_Invalid_Stream_Read:
		return "FT_ERROR_Invalid_Stream_Read"
	case FT_ERROR_Invalid_Stream_Operation:
		return "FT_ERROR_Invalid_Stream_Operation"
	case FT_ERROR_Invalid_Frame_Operation:
		return "FT_ERROR_Invalid_Frame_Operation"
	case FT_ERROR_Nested_Frame_Access:
		return "FT_ERROR_Nested_Frame_Access"
	case FT_ERROR_Invalid_Frame_Read:
		return "FT_ERROR_Invalid_Frame_Read"
	case FT_ERROR_Raster_Uninitialized:
		return "FT_ERROR_Raster_Uninitialized"
	case FT_ERROR_Raster_Corrupted:
		return "FT_ERROR_Raster_Corrupted"
	case FT_ERROR_Raster_Overflow:
		return "FT_ERROR_Raster_Overflow"
	case FT_ERROR_Raster_Negative_Height:
		return "FT_ERROR_Raster_Negative_Height"
	case FT_ERROR_Too_Many_Caches:
		return "FT_ERROR_Too_Many_Caches"
	case FT_ERROR_Invalid_Opcode:
		return "FT_ERROR_Invalid_Opcode"
	case FT_ERROR_Too_Few_Arguments:
		return "FT_ERROR_Too_Few_Arguments"
	case FT_ERROR_Stack_Overflow:
		return "FT_ERROR_Stack_Overflow"
	case FT_ERROR_Code_Overflow:
		return "FT_ERROR_Code_Overflow"
	case FT_ERROR_Bad_Argument:
		return "FT_ERROR_Bad_Argument"
	case FT_ERROR_Divide_By_Zero:
		return "FT_ERROR_Divide_By_Zero"
	case FT_ERROR_Invalid_Reference:
		return "FT_ERROR_Invalid_Reference"
	case FT_ERROR_Debug_OpCode:
		return "FT_ERROR_Debug_OpCode"
	case FT_ERROR_ENDF_In_Exec_Stream:
		return "FT_ERROR_ENDF_In_Exec_Stream"
	case FT_ERROR_Nested_DEFS:
		return "FT_ERROR_Nested_DEFS"
	case FT_ERROR_Invalid_CodeRange:
		return "FT_ERROR_Invalid_CodeRange"
	case FT_ERROR_Execution_Too_Long:
		return "FT_ERROR_Execution_Too_Long"
	case FT_ERROR_Too_Many_Function_Defs:
		return "FT_ERROR_Too_Many_Function_Defs"
	case FT_ERROR_Too_Many_Instruction_Defs:
		return "FT_ERROR_Too_Many_Instruction_Defs"
	case FT_ERROR_Table_Missing:
		return "FT_ERROR_Table_Missing"
	case FT_ERROR_Horiz_Header_Missing:
		return "FT_ERROR_Horiz_Header_Missing"
	case FT_ERROR_Locations_Missing:
		return "FT_ERROR_Locations_Missing"
	case FT_ERROR_Name_Table_Missing:
		return "FT_ERROR_Name_Table_Missing"
	case FT_ERROR_CMap_Table_Missing:
		return "FT_ERROR_CMap_Table_Missing"
	case FT_ERROR_Hmtx_Table_Missing:
		return "FT_ERROR_Hmtx_Table_Missing"
	case FT_ERROR_Post_Table_Missing:
		return "FT_ERROR_Post_Table_Missing"
	case FT_ERROR_Invalid_Horiz_Metrics:
		return "FT_ERROR_Invalid_Horiz_Metrics"
	case FT_ERROR_Invalid_CharMap_Format:
		return "FT_ERROR_Invalid_CharMap_Format"
	case FT_ERROR_Invalid_PPem:
		return "FT_ERROR_Invalid_PPem"
	case FT_ERROR_Invalid_Vert_Metrics:
		return "FT_ERROR_Invalid_Vert_Metrics"
	case FT_ERROR_Could_Not_Find_Context:
		return "FT_ERROR_Could_Not_Find_Context"
	case FT_ERROR_Invalid_Post_Table_Format:
		return "FT_ERROR_Invalid_Post_Table_Format"
	case FT_ERROR_Invalid_Post_Table:
		return "FT_ERROR_Invalid_Post_Table"
	case FT_ERROR_Syntax_Error:
		return "FT_ERROR_Syntax_Error"
	case FT_ERROR_Stack_Underflow:
		return "FT_ERROR_Stack_Underflow"
	case FT_ERROR_Ignore:
		return "FT_ERROR_Ignore"
	case FT_ERROR_No_Unicode_Glyph_Name:
		return "FT_ERROR_No_Unicode_Glyph_Name"
	case FT_ERROR_Glyph_Too_Big:
		return "FT_ERROR_Glyph_Too_Big"
	case FT_ERROR_Missing_Startfont_Field:
		return "FT_ERROR_Missing_Startfont_Field"
	case FT_ERROR_Missing_Font_Field:
		return "FT_ERROR_Missing_Font_Field"
	case FT_ERROR_Missing_Size_Field:
		return "FT_ERROR_Missing_Size_Field"
	case FT_ERROR_Missing_Fontboundingbox_Field:
		return "FT_ERROR_Missing_Fontboundingbox_Field"
	case FT_ERROR_Missing_Chars_Field:
		return "FT_ERROR_Missing_Chars_Field"
	case FT_ERROR_Missing_Startchar_Field:
		return "FT_ERROR_Missing_Startchar_Field"
	case FT_ERROR_Missing_Encoding_Field:
		return "FT_ERROR_Missing_Encoding_Field"
	case FT_ERROR_Missing_Bbx_Field:
		return "FT_ERROR_Missing_Bbx_Field"
	case FT_ERROR_Bbx_Too_Big:
		return "FT_ERROR_Bbx_Too_Big"
	case FT_ERROR_Corrupted_Font_Header:
		return "FT_ERROR_Corrupted_Font_Header"
	case FT_ERROR_Corrupted_Font_Glyphs:
		return "FT_ERROR_Corrupted_Font_Glyphs"
	default:
		return "[?? Invalid FT_ERROR value]" // TODO: Add error codes
	}
}
