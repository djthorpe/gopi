/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

/*
  #cgo CFLAGS:   -I/usr/include/freetype2 -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lfreetype -lOpenVG
  #include <ft2build.h>
  #include <freetype.h>
  #include <VG/openvg.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Configuration when creating the OpenVGFont driver
type VGFont struct {
	// Pixels per inch, the density of pixels on the screen
	// If you set this to zero then pixels are used instead when
	// requesting sizes when drawing
	PPI uint
}

// VGFont driver
type vgfDriver struct {
	handle C.FT_Library
	log    gopi.Logger
	count  uint
	lock   sync.Mutex
	faces  []*vgfFace
	ppi    uint
}

// vgfFace represents a loaded TTF face
type vgfFace struct {
	count  uint
	handle C.FT_Face
	font   C.VGFont
	path   string
	ppi    uint
}

// vgfEncoding represents charmap encoding
type vgfEncoding string

// glyph bitmap format
type VGFontBitmapPixelMode C.FT_Pixel_Mode

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_FONT_NONE                 C.VGFont    = C.VGFont(0)
	VG_FONT_PPI_DEFAULT          uint        = 72
	VG_FONT_FT_ERROR_NONE        C.FT_Error  = C.FT_Error(0)
	VG_FONT_FT_ENCODING_UNICODE  vgfEncoding = "unic"
	VG_FONT_FT_STYLE_FLAG_ITALIC C.FT_Long   = (1 << 0)
	VG_FONT_FT_STYLE_FLAG_BOLD   C.FT_Long   = (1 << 1)
)

const (
	VG_FONT_FT_LOAD_DEFAULT                     uint32 = 0
	VG_FONT_FT_LOAD_NO_SCALE                    uint32 = (1 << 0)
	VG_FONT_FT_LOAD_NO_HINTING                  uint32 = (1 << 1)
	VG_FONT_FT_LOAD_RENDER                      uint32 = (1 << 2)
	VG_FONT_FT_LOAD_NO_BITMAP                   uint32 = (1 << 3)
	VG_FONT_FT_LOAD_VERTICAL_LAYOUT             uint32 = (1 << 4)
	VG_FONT_FT_LOAD_FORCE_AUTOHINT              uint32 = (1 << 5)
	VG_FONT_FT_LOAD_CROP_BITMAP                 uint32 = (1 << 6)
	VG_FONT_FT_LOAD_PEDANTIC                    uint32 = (1 << 7)
	VG_FONT_FT_LOAD_IGNORE_GLOBAL_ADVANCE_WIDTH uint32 = (1 << 9)
	VG_FONT_FT_LOAD_NO_RECURSE                  uint32 = (1 << 10)
	VG_FONT_FT_LOAD_IGNORE_TRANSFORM            uint32 = (1 << 11)
	VG_FONT_FT_LOAD_MONOCHROME                  uint32 = (1 << 12)
	VG_FONT_FT_LOAD_LINEAR_DESIGN               uint32 = (1 << 13)
	VG_FONT_FT_LOAD_NO_AUTOHINT                 uint32 = (1 << 15)
	VG_FONT_FT_LOAD_COLOR                       uint32 = (1 << 20)
	VG_FONT_FT_LOAD_COMPUTE_METRICS             uint32 = (1 << 21)
)

const (
	VG_FONT_FT_PIXEL_MODE_NONE  VGFontBitmapPixelMode = iota
	VG_FONT_FT_PIXEL_MODE_MONO                        // 1 bit per pixel (unsupported)
	VG_FONT_FT_PIXEL_MODE_GRAY                        // 8 bits per pixel
	VG_FONT_FT_PIXEL_MODE_GRAY2                       // 2 bits per pixel (unsupported)
	VG_FONT_FT_PIXEL_MODE_GRAY4                       // 4 bits per pixel (unsupported)
	VG_FONT_FT_PIXEL_MODE_LCD                         // 8 bits per pixel, horizontal LCD (unsupported)
	VG_FONT_FT_PIXEL_MODE_LCD_V                       // 8 bits per pixel, vertical LCD (unsupported)
	VG_FONT_FT_PIXEL_MODE_BGRA                        // 32 bits per pixel (unsupported)
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

// Open the driver
func (config VGFont) Open(log gopi.Logger) (gopi.Driver, error) {
	this := new(vgfDriver)
	this.log = log
	this.count = 0
	this.ppi = config.PPI

	this.log.Debug2("<rpi.VGFont>Open")

	// initialise
	this.lock.Lock()
	defer this.lock.Unlock()
	if err := this.vgfontInit(); err != nil {
		return nil, err
	}

	// Set PPI to default value if it's not set. A value of 0 should
	// probably be used so that pixels are used instead of points, but we'll
	// leave that enhancement to the future
	if this.ppi == 0 {
		this.ppi = VG_FONT_PPI_DEFAULT
	}

	// Create faces structure
	this.faces = make([]*vgfFace, 0)

	// Success
	return this, nil
}

// Close the driver
func (this *vgfDriver) Close() error {
	this.log.Debug2("<rpi.VGFont>Close")

	// destroy faces
	for {
		if len(this.faces) == 0 {
			break
		}
		this.DestroyFace(this.faces[0])
	}

	// Destroy
	this.lock.Lock()
	defer this.lock.Unlock()
	if err := this.vgfontDestroy(); err != nil {
		return err
	}

	return nil
}

// Return human-readable form of driver
func (this *vgfDriver) String() string {
	return fmt.Sprintf("<rpi.VGFont>{ handle=%v ppi=%v faces=%v }", this.handle, this.ppi, this.faces)
}

// Return string array of families
func (this *vgfDriver) GetFamilies() []string {
	family_map := make(map[string]bool, 0)
	family_array := make([]string, 0)
	for _, face := range this.faces {
		family := face.GetFamily()
		if _, exists := family_map[family]; exists {
			continue
		}
		family_map[family] = true
		family_array = append(family_array, family)
	}
	return family_array
}

// Return faces in a family and/or with a particular set of attributes
func (this *vgfDriver) GetFaces(family string, flags khronos.VGFontStyleFlags) []khronos.VGFace {
	faces := make([]khronos.VGFace, 0)
	for _, face := range this.faces {
		if family != "" && family != face.GetFamily() {
			continue
		}
		switch flags {
		case khronos.VG_FONT_STYLE_ANY:
			faces = append(faces, face)
			break
		case khronos.VG_FONT_STYLE_REGULAR:
			if face.IsBold() == false && face.IsItalic() == false {
				faces = append(faces, face)
			}
			break
		case khronos.VG_FONT_STYLE_BOLD:
			if face.IsBold() == true && face.IsItalic() == false {
				faces = append(faces, face)
			}
			break
		case khronos.VG_FONT_STYLE_ITALIC:
			if face.IsBold() == false && face.IsItalic() == true {
				faces = append(faces, face)
			}
			break
		case khronos.VG_FONT_STYLE_BOLDITALIC:
			if face.IsBold() == true && face.IsItalic() == true {
				faces = append(faces, face)
			}
			break
		}
	}
	return faces
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS: Open and Destroy Faces

func (this *vgfDriver) OpenFace(path string) (khronos.VGFace, error) {
	return this.OpenFaceAtIndex(path, 0)
}

func (this *vgfDriver) OpenFaceAtIndex(path string, index uint) (khronos.VGFace, error) {
	this.log.Debug2("<rpi.VGFace>OpenFaceAtIndex path=%v index=%v", path, index)

	var err error
	face := new(vgfFace)
	face.path = path
	face.count = this.count
	face.ppi = this.ppi

	this.lock.Lock()
	defer this.lock.Unlock()
	this.count += 1
	face.handle, err = this.vgfontLoadFace(path, index)
	if err != nil {
		return nil, err
	}

	// Set Unicode
	if err := this.vgfontSelectCharmap(face.handle, VG_FONT_FT_ENCODING_UNICODE); err != nil {
		this.vgfontDoneFace(face.handle)
		return nil, err
	}

	// VG Create Font
	//face.font = C.vgCreateFont(C.VGint(face.GetNumGlyphs()))
	//if face.font == VG_FONT_NONE {
	//	this.vgfontDoneFace(face.handle)
	//	return nil, vgGetError(vgErrorType(C.vgGetError()))
	//}

	// Load Glyphs
	//if err := this.LoadGlyphs(face, 64.0, 0.0); err != nil {
	//	this.vgfontDoneFace(face.handle)
	//	C.vgDestroyFont(face.font)
	//	return nil, err
	//}

	// Add face to list of faces
	this.faces = append(this.faces, face)

	return face, nil
}

func (this *vgfDriver) OpenFacesAtPath(path string, callback func(path string, info os.FileInfo) bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if callback(path, info) == false {
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}
		if info.IsDir() {
			return nil
		}

		// Open zero-indexed face
		face, err := this.OpenFace(path)
		if err != nil {
			return err
		}

		// If there are more faces in the file, then load these too
		if face.GetNumFaces() > uint(1) {
			for i := uint(1); i < face.GetNumFaces(); i++ {
				_, err := this.OpenFaceAtIndex(path, i)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	return err
}

func (this *vgfDriver) DestroyFace(face khronos.VGFace) error {
	this.log.Debug2("<rpi.VGFace>DestroyFace %v", face)

	// Remove the face from the list of faces
	j := -1
	for i, existing_face := range this.faces {
		if existing_face.count == face.(*vgfFace).count {
			j = i
			break
		}
	}
	if j >= 0 {
		this.faces = append(this.faces[:j], this.faces[j+1:]...)
	}

	// Destroy the VGFont
	C.vgDestroyFont(face.(*vgfFace).font)

	// Destroy the face
	if err := this.vgfontDoneFace(face.(*vgfFace).GetHandle()); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS: Face information

func (this *vgfFace) String() string {
	return fmt.Sprintf("<rpi.VGFace>{ id=%v name=%v index=%v family=%v style=%v is_bold=%v is_italic=%v num_faces=%v num_glyphs=%v }", this.count, this.GetName(), this.GetIndex(), this.GetFamily(), this.GetStyle(), this.IsBold(), this.IsItalic(), this.GetNumFaces(), this.GetNumGlyphs())
}

func (this *vgfFace) GetName() string {
	return path.Base(this.path)
}

func (this *vgfFace) GetHandle() C.FT_Face {
	return this.handle
}

func (this *vgfFace) GetFamily() string {
	return C.GoString((*C.char)(this.handle.family_name))
}

func (this *vgfFace) GetStyle() string {
	return C.GoString((*C.char)(this.handle.style_name))
}

func (this *vgfFace) GetIndex() uint {
	return uint(this.handle.face_index)
}

func (this *vgfFace) GetNumFaces() uint {
	return uint(this.handle.num_faces)
}

func (this *vgfFace) GetNumGlyphs() uint {
	return uint(this.handle.num_glyphs)
}

func (this *vgfFace) IsBold() bool {
	return (this.handle.style_flags & VG_FONT_FT_STYLE_FLAG_BOLD) != 0
}

func (this *vgfFace) IsItalic() bool {
	return (this.handle.style_flags & VG_FONT_FT_STYLE_FLAG_ITALIC) != 0
}

////////////////////////////////////////////////////////////////////////////////
// BITMAP FUNCTIONS

func (this *vgfFace) SetSize(points float32) error {
	if this.ppi == 0 {
		// Set as pixels
		if ret := C.FT_Set_Pixel_Sizes(this.handle, 0, C.FT_UInt(points)); ret != VG_FONT_FT_ERROR_NONE {
			return vgfontGetError(ret)
		}
	} else {
		// Set as points
		if ret := C.FT_Set_Char_Size(this.handle, 0, C.FT_F26Dot6(points*64.0), 0, C.FT_UInt(this.ppi)); ret != VG_FONT_FT_ERROR_NONE {
			return vgfontGetError(ret)
		}
	}
	return nil
}

// This method returns a bitmap for a rune. The returned values are a pointer
// to the bitmap pixels, the width and height of the bitmap data, the advancement
// required to draw the next bitmap rune, the stride value for the bitmap (number
// of pixels per row) and an error condition if the rune could not be found in
// the font face (for example) or the size had not yet been set.
func (this *vgfFace) LoadBitmapForRune(value rune) (uintptr, VGFontBitmapPixelMode, khronos.EGLSize, khronos.EGLSize, uint, error) {

	// Get Glyph
	glyph_index := C.FT_Get_Char_Index(this.handle, C.FT_ULong(value))
	if glyph_index == 0 {
		return uintptr(0), VG_FONT_FT_PIXEL_MODE_NONE, khronos.EGLZeroSize, khronos.EGLZeroSize, 0, ErrFontInvalidCharacterCode
	}

	// Render Glyph
	ret := C.FT_Load_Glyph(this.handle, glyph_index, C.FT_Int32(VG_FONT_FT_LOAD_RENDER))
	if ret != VG_FONT_FT_ERROR_NONE {
		return uintptr(0), VG_FONT_FT_PIXEL_MODE_NONE, khronos.EGLZeroSize, khronos.EGLZeroSize, 0, vgfontGetError(ret)
	}

	// Compute relevant information
	bitmap := this.handle.glyph.bitmap
	pixel_mode := VGFontBitmapPixelMode(bitmap.pixel_mode)
	size := khronos.EGLSize{Width: uint(bitmap.width), Height: uint(bitmap.rows)}
	advance := khronos.EGLSize{Width: uint(this.handle.glyph.advance.x >> 6), Height: uint(this.handle.glyph.advance.y >> 6)}
	stride := uint(bitmap.pitch)

	// Success
	return uintptr(unsafe.Pointer(bitmap.buffer)), pixel_mode, size, advance, stride, nil
}

////////////////////////////////////////////////////////////////////////////////
// CONVERT GLYPHS AT PARTICULAR SIZE

// Load Glyphs at (w,h) point size
func (this *vgfDriver) LoadGlyphs(face *vgfFace, w, h float32) error {
	ret := C.FT_Set_Char_Size(face.handle, C.FT_F26Dot6(w*64.0), C.FT_F26Dot6(h*64.0), C.FT_UInt(this.ppi), C.FT_UInt(this.ppi))
	if ret != VG_FONT_FT_ERROR_NONE {
		return vgfontGetError(ret)
	}

	var glyph C.FT_UInt
	handle := C.FT_Get_First_Char(face.handle, (*C.FT_UInt)(unsafe.Pointer(&glyph)))
	for glyph != C.FT_UInt(0) {

		// Load Glyph into font by converting into paths
		ret := C.FT_Load_Glyph(face.handle, glyph, C.FT_Int32(VG_FONT_FT_LOAD_DEFAULT))
		if ret != VG_FONT_FT_ERROR_NONE {
			return vgfontGetError(ret)
		}
		if err := this.vgfontLoadGlyphToFont(face, glyph); err != nil {
			return err
		}

		// Move to next glyph
		handle = C.FT_Get_Next_Char(face.handle, handle, (*C.FT_UInt)(unsafe.Pointer(&glyph)))
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS

func (this *vgfDriver) vgfontInit() error {
	return vgfontGetError(C.FT_Init_FreeType(&this.handle))
}

func (this *vgfDriver) vgfontDestroy() error {
	return vgfontGetError(C.FT_Done_FreeType(this.handle))
}

func (this *vgfDriver) vgfontLoadFace(path string, index uint) (C.FT_Face, error) {
	var face C.FT_Face
	ret := C.FT_New_Face(this.handle, C.CString(path), C.FT_Long(index), &face)
	return face, vgfontGetError(ret)
}

func (this *vgfDriver) vgfontDoneFace(handle C.FT_Face) error {
	return vgfontGetError(C.FT_Done_Face(handle))
}

func (this *vgfDriver) vgfontSelectCharmap(handle C.FT_Face, encoding vgfEncoding) error {
	code := C.FT_Encoding(uint32(encoding[3]) | uint32(encoding[2])<<8 | uint32(encoding[1])<<16 | uint32(encoding[0])<<24)
	return vgfontGetError(C.FT_Select_Charmap(handle, code))
}

func (this *vgfDriver) vgfontLoadGlyphToFont(face *vgfFace, glyph C.FT_UInt) error {
	// create a path
	path := VG_PATH_HANDLE_NONE
	outline := face.handle.glyph.outline
	if outline.n_contours != 0 {
		path = C.vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_F, C.VGfloat(1.0), C.VGfloat(0.0), C.VGint(0), C.VGint(0), C.VGbitfield(VG_PATH_CAPABILITY_ALL))
		if path == VG_PATH_HANDLE_NONE {
			return vgGetError()
		}
	}

	/* TODO
	   	origin := []C.VGfloat{ C.VGfloat(0), C.VGfloat(0) }
	   	escapement :=
	         VGfloat escapement[] = { float_from_26_6(font->ft_face->glyph->advance.x), float_from_26_6(font->ft_face->glyph->advance.y) };
	         vgSetGlyphToPath(font->vg_font, glyph_index, vg_path, VG_FALSE, origin, escapement);
	*/

	// destroy path
	if path != VG_PATH_HANDLE_NONE {
		C.vgDestroyPath(path)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m VGFontBitmapPixelMode) String() string {
	switch m {
	case VG_FONT_FT_PIXEL_MODE_NONE:
		return "VG_FONT_FT_PIXEL_MODE_NONE"
	case VG_FONT_FT_PIXEL_MODE_MONO:
		return "VG_FONT_FT_PIXEL_MODE_MONO"
	case VG_FONT_FT_PIXEL_MODE_GRAY:
		return "VG_FONT_FT_PIXEL_MODE_GRAY"
	case VG_FONT_FT_PIXEL_MODE_GRAY2:
		return "VG_FONT_FT_PIXEL_MODE_GRAY2"
	case VG_FONT_FT_PIXEL_MODE_GRAY4:
		return "VG_FONT_FT_PIXEL_MODE_GRAY4"
	case VG_FONT_FT_PIXEL_MODE_LCD:
		return "VG_FONT_FT_PIXEL_MODE_LCD"
	case VG_FONT_FT_PIXEL_MODE_LCD_V:
		return "VG_FONT_FT_PIXEL_MODE_LCD_V"
	case VG_FONT_FT_PIXEL_MODE_BGRA:
		return "VG_FONT_FT_PIXEL_MODE_BGRA"
	default:
		return "[?? Invalid VGFontBitmapPixelMode value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// ERROR HANDLING

var (
	/* generic errors */
	ErrFontCannotOpenResource   = errors.New("Cannot_Open_Resource")  // 0x01
	ErrFontUnknownFileFormat    = errors.New("Unknown_File_Format")   // 0x02
	ErrFontInvalidFileFormat    = errors.New("Invalid_File_Format")   // 0x03
	ErrFontInvalidVersion       = errors.New("Invalid_Version")       // 0x04
	ErrFontLowerModuleVersion   = errors.New("Lower_Module_Version")  // 0x05
	ErrFontInvalidArgument      = errors.New("Invalid_Argument")      // 0x06
	ErrFontUnimplementedFeature = errors.New("Unimplemented_Feature") // 0x07
	ErrFontInvalidTable         = errors.New("Invalid_Table")         // 0x08
	ErrFontInvalidOffset        = errors.New("Invalid_Offset")        // 0x09
	ErrFontArrayTooLarge        = errors.New("Array_Too_Large")       // 0x0A
	ErrFontMissingModule        = errors.New("Missing_Module")        // 0x0B
	ErrFontMissingProperty      = errors.New("Missing_Property")      // 0x0C

	/* glyph/character errors */
	ErrFontInvalidGlyphIndex    = errors.New("Invalid_Glyph_Index")    // 0x10
	ErrFontInvalidCharacterCode = errors.New("Invalid_Character_Code") // 0x11
	ErrFontInvalidGlyphFormat   = errors.New("Invalid_Glyph_Format")   // 0x12
	ErrFontCannotRenderGlyph    = errors.New("Cannot_Render_Glyph")    // 0x13
	ErrFontInvalidOutline       = errors.New("Invalid_Outline")        // 0x14
	ErrFontInvalidComposite     = errors.New("Invalid_Composite")      // 0x15
	ErrFontTooManyHints         = errors.New("Too_Many_Hints")         // 0x16
	ErrFontInvalidPixelSize     = errors.New("Invalid_Pixel_Size")     // 0x17

	/* handle errors */
	ErrFontInvalidHandle        = errors.New("Invalid_Handle")         // 0x20
	ErrFontInvalidLibraryHandle = errors.New("Invalid_Library_Handle") // 0x21
	ErrFontInvalidDriverHandle  = errors.New("Invalid_Driver_Handle")  // 0x22
	ErrFontInvalidFaceHandle    = errors.New("Invalid_Face_Handle")    // 0x23
	ErrFontInvalidSizeHandle    = errors.New("Invalid_Size_Handle")    // 0x24
	ErrFontInvalidSlotHandle    = errors.New("Invalid_Slot_Handle")    // 0x25
	ErrFontInvalidCharMapHandle = errors.New("Invalid_CharMap_Handle") // 0x26
	ErrFontInvalidCacheHandle   = errors.New("Invalid_Cache_Handle")   // 0x27
	ErrFontInvalidStreamHandle  = errors.New("Invalid_Stream_Handle")  // 0x28

	/* driver errors */
	ErrFontTooManyDrivers    = errors.New("Too_Many_Drivers")    // 0x30
	ErrFontTooManyExtensions = errors.New("Too_Many_Extensions") // 0x31

	/* memory errors */
	ErrFontOutOfMemory    = errors.New("Out_Of_Memory")   // 0x40
	ErrFontUnlistedObject = errors.New("Unlisted_Object") // 0x41

	/* stream errors */
	ErrFontCannotOpenStream       = errors.New("Cannot_Open_Stream")       // 0x51
	ErrFontInvalidStreamSeek      = errors.New("Invalid_Stream_Seek")      // 0x52
	ErrFontInvalidStreamSkip      = errors.New("Invalid_Stream_Skip")      // 0x53
	ErrFontInvalidStreamRead      = errors.New("Invalid_Stream_Read")      // 0x54
	ErrFontInvalidStreamOperation = errors.New("Invalid_Stream_Operation") // 0x55
	ErrFontInvalidFrameOperation  = errors.New("Invalid_Frame_Operation")  // 0x56
	ErrFontNestedFrameAccess      = errors.New("Nested_Frame_Access")      // 0x57
	ErrFontInvalidFrameRead       = errors.New("Invalid_Frame_Read")       // 0x58

	/* raster errors */
	ErrFontRasterUninitialized  = errors.New("Raster_Uninitialized")   // 0x60
	ErrFontRasterCorrupted      = errors.New("Raster_Corrupted")       // 0x61
	ErrFontRasterOverflow       = errors.New("Raster_Overflow")        // 0x62
	ErrFontRasterNegativeHeight = errors.New("Raster_Negative_Height") // 0x63

	/* cache errors */
	ErrFontTooManyCaches = errors.New("Too_Many_Caches") // 0x70

	/* TrueType and SFNT errors */
	ErrFontInvalidOpcode          = errors.New("Invalid_Opcode")            // 0x80
	ErrFontTooFewArguments        = errors.New("Too_Few_Arguments")         // 0x81
	ErrFontStackOverflow          = errors.New("Stack_Overflow")            // 0x82
	ErrFontCodeOverflow           = errors.New("Code_Overflow")             // 0x83
	ErrFontBadArgument            = errors.New("Bad_Argument")              // 0x84
	ErrFontDivideByZero           = errors.New("Divide_By_Zero")            // 0x85
	ErrFontInvalidReference       = errors.New("Invalid_Reference")         // 0x86
	ErrFontDebugOpCode            = errors.New("Debug_OpCode")              // 0x87
	ErrFontENDFInExecStream       = errors.New("ENDF_In_Exec_Stream")       // 0x88
	ErrFontNestedDEFS             = errors.New("Nested_DEFS")               // 0x89
	ErrFontInvalidCodeRange       = errors.New("Invalid_CodeRange")         // 0x8A
	ErrFontExecutionTooLong       = errors.New("Execution_Too_Long")        // 0x8B
	ErrFontTooManyFunctionDefs    = errors.New("Too_Many_Function_Defs")    // 0x8C
	ErrFontTooManyInstructionDefs = errors.New("Too_Many_Instruction_Defs") // 0x8D
	ErrFontTableMissing           = errors.New("Table_Missing")             // 0x8E
	ErrFontHorizHeaderMissing     = errors.New("Horiz_Header_Missing")      // 0x8F
	ErrFontLocationsMissing       = errors.New("Locations_Missing")         // 0x90
	ErrFontNameTableMissing       = errors.New("Name_Table_Missing")        // 0x91
	ErrFontCMapTableMissing       = errors.New("CMap_Table_Missing")        // 0x92
	ErrFontHmtxTableMissing       = errors.New("Hmtx_Table_Missing")        // 0x93
	ErrFontPostTableMissing       = errors.New("Post_Table_Missing")        // 0x94
	ErrFontInvalidHorizMetrics    = errors.New("Invalid_Horiz_Metrics")     // 0x95
	ErrFontInvalidCharMapFormat   = errors.New("Invalid_CharMap_Format")    // 0x96
	ErrFontInvalidPPem            = errors.New("Invalid_PPem")              // 0x97
	ErrFontInvalidVertMetrics     = errors.New("Invalid_Vert_Metrics")      // 0x98
	ErrFontCouldNotFindContext    = errors.New("Could_Not_Find_Context")    // 0x99
	ErrFontInvalidPostTableFormat = errors.New("Invalid_Post_Table_Format") // 0x9A
	ErrFontInvalidPostTable       = errors.New("Invalid_Post_Table")        // 0x9B

	/* CFF, CID, and Type 1 errors */
	ErrFontSyntaxError        = errors.New("Syntax_Error")          // 0xA0
	ErrFontStackUnderflow     = errors.New("Stack_Underflow")       // 0xA1
	ErrFontIgnore             = errors.New("Ignore")                // 0xA2
	ErrFontNoUnicodeGlyphName = errors.New("No_Unicode_Glyph_Name") // 0xA3
	ErrFontGlyphTooBig        = errors.New("Glyph_Too_Big")         // 0xA4

	/* BDF errors */
	ErrFontMissingStartfontField       = errors.New("Missing_Startfont_Field")       // 0xB0
	ErrFontMissingFontField            = errors.New("Missing_Font_Field")            // 0xB1
	ErrFontMissingSizeField            = errors.New("Missing_Size_Field")            // 0xB2
	ErrFontMissingFontboundingboxField = errors.New("Missing_Fontboundingbox_Field") // 0xB3
	ErrFontMissingCharsField           = errors.New("Missing_Chars_Field")           // 0xB4
	ErrFontMissingStartcharField       = errors.New("Missing_Startchar_Field")       // 0xB5
	ErrFontMissingEncodingField        = errors.New("Missing_Encoding_Field")        // 0xB6
	ErrFontMissingBbxField             = errors.New("Missing_Bbx_Field")             // 0xB7
	ErrFontBbxTooBig                   = errors.New("Bbx_Too_Big")                   // 0xB8
	ErrFontCorruptedFontHeader         = errors.New("Corrupted_Font_Header")         // 0xB9
	ErrFontCorruptedFontGlyphs         = errors.New("Corrupted_Font_Glyphs")         // 0xBA
)

func vgfontGetError(code C.FT_Error) error {
	if code == C.FT_Error(0) {
		return nil
	}
	switch code {
	case C.FT_Error(0x01):
		return ErrFontCannotOpenResource
	case C.FT_Error(0x02):
		return ErrFontUnknownFileFormat
	case C.FT_Error(0x03):
		return ErrFontInvalidFileFormat
	case C.FT_Error(0x04):
		return ErrFontInvalidVersion
	case C.FT_Error(0x05):
		return ErrFontLowerModuleVersion
	case C.FT_Error(0x06):
		return ErrFontInvalidArgument
	case C.FT_Error(0x07):
		return ErrFontUnimplementedFeature
	case C.FT_Error(0x08):
		return ErrFontInvalidTable
	case C.FT_Error(0x09):
		return ErrFontInvalidOffset
	case C.FT_Error(0x0A):
		return ErrFontArrayTooLarge
	case C.FT_Error(0x0B):
		return ErrFontMissingModule
	case C.FT_Error(0x0C):
		return ErrFontMissingProperty

	/* glyph/character errors */
	case C.FT_Error(0x10):
		return ErrFontInvalidGlyphIndex
	case C.FT_Error(0x11):
		return ErrFontInvalidCharacterCode
	case C.FT_Error(0x12):
		return ErrFontInvalidGlyphFormat
	case C.FT_Error(0x13):
		return ErrFontCannotRenderGlyph
	case C.FT_Error(0x14):
		return ErrFontInvalidOutline
	case C.FT_Error(0x15):
		return ErrFontInvalidComposite
	case C.FT_Error(0x16):
		return ErrFontTooManyHints
	case C.FT_Error(0x17):
		return ErrFontInvalidPixelSize

	/* handle errors */
	case C.FT_Error(0x20):
		return ErrFontInvalidHandle
	case C.FT_Error(0x21):
		return ErrFontInvalidLibraryHandle
	case C.FT_Error(0x22):
		return ErrFontInvalidDriverHandle
	case C.FT_Error(0x23):
		return ErrFontInvalidFaceHandle
	case C.FT_Error(0x24):
		return ErrFontInvalidSizeHandle
	case C.FT_Error(0x25):
		return ErrFontInvalidSlotHandle
	case C.FT_Error(0x26):
		return ErrFontInvalidCharMapHandle
	case C.FT_Error(0x27):
		return ErrFontInvalidCacheHandle
	case C.FT_Error(0x28):
		return ErrFontInvalidStreamHandle

	/* driver errors */
	case C.FT_Error(0x30):
		return ErrFontTooManyDrivers
	case C.FT_Error(0x31):
		return ErrFontTooManyExtensions

	/* memory errors */
	case C.FT_Error(0x40):
		return ErrFontOutOfMemory
	case C.FT_Error(0x41):
		return ErrFontUnlistedObject

	/* stream errors */
	case C.FT_Error(0x51):
		return ErrFontCannotOpenStream
	case C.FT_Error(0x52):
		return ErrFontInvalidStreamSeek
	case C.FT_Error(0x53):
		return ErrFontInvalidStreamSkip
	case C.FT_Error(0x54):
		return ErrFontInvalidStreamRead
	case C.FT_Error(0x55):
		return ErrFontInvalidStreamOperation
	case C.FT_Error(0x56):
		return ErrFontInvalidFrameOperation
	case C.FT_Error(0x57):
		return ErrFontNestedFrameAccess
	case C.FT_Error(0x58):
		return ErrFontInvalidFrameRead

	/* raster errors */
	case C.FT_Error(0x60):
		return ErrFontRasterUninitialized
	case C.FT_Error(0x61):
		return ErrFontRasterCorrupted
	case C.FT_Error(0x62):
		return ErrFontRasterOverflow
	case C.FT_Error(0x63):
		return ErrFontRasterNegativeHeight

	/* cache errors */
	case C.FT_Error(0x70):
		return ErrFontTooManyCaches

	/* TrueType and SFNT errors */
	case C.FT_Error(0x80):
		return ErrFontInvalidOpcode
	case C.FT_Error(0x81):
		return ErrFontTooFewArguments
	case C.FT_Error(0x82):
		return ErrFontStackOverflow
	case C.FT_Error(0x83):
		return ErrFontCodeOverflow
	case C.FT_Error(0x84):
		return ErrFontBadArgument
	case C.FT_Error(0x85):
		return ErrFontDivideByZero
	case C.FT_Error(0x86):
		return ErrFontInvalidReference
	case C.FT_Error(0x87):
		return ErrFontDebugOpCode
	case C.FT_Error(0x88):
		return ErrFontENDFInExecStream
	case C.FT_Error(0x89):
		return ErrFontNestedDEFS
	case C.FT_Error(0x8A):
		return ErrFontInvalidCodeRange
	case C.FT_Error(0x8B):
		return ErrFontExecutionTooLong
	case C.FT_Error(0x8C):
		return ErrFontTooManyFunctionDefs
	case C.FT_Error(0x8D):
		return ErrFontTooManyInstructionDefs
	case C.FT_Error(0x8E):
		return ErrFontTableMissing
	case C.FT_Error(0x8F):
		return ErrFontHorizHeaderMissing
	case C.FT_Error(0x90):
		return ErrFontLocationsMissing
	case C.FT_Error(0x91):
		return ErrFontNameTableMissing
	case C.FT_Error(0x92):
		return ErrFontCMapTableMissing
	case C.FT_Error(0x93):
		return ErrFontHmtxTableMissing
	case C.FT_Error(0x94):
		return ErrFontPostTableMissing
	case C.FT_Error(0x95):
		return ErrFontInvalidHorizMetrics
	case C.FT_Error(0x96):
		return ErrFontInvalidCharMapFormat
	case C.FT_Error(0x97):
		return ErrFontInvalidPPem
	case C.FT_Error(0x98):
		return ErrFontInvalidVertMetrics
	case C.FT_Error(0x99):
		return ErrFontCouldNotFindContext
	case C.FT_Error(0x9A):
		return ErrFontInvalidPostTableFormat
	case C.FT_Error(0x9B):
		return ErrFontInvalidPostTable

	/* CFF, CID, and Type 1 errors */
	case C.FT_Error(0xA0):
		return ErrFontSyntaxError
	case C.FT_Error(0xA1):
		return ErrFontStackUnderflow
	case C.FT_Error(0xA2):
		return ErrFontIgnore
	case C.FT_Error(0xA3):
		return ErrFontNoUnicodeGlyphName
	case C.FT_Error(0xA4):
		return ErrFontGlyphTooBig

	/* BDF errors */
	case C.FT_Error(0xB0):
		return ErrFontMissingStartfontField
	case C.FT_Error(0xB1):
		return ErrFontMissingFontField
	case C.FT_Error(0xB2):
		return ErrFontMissingSizeField
	case C.FT_Error(0xB3):
		return ErrFontMissingFontboundingboxField
	case C.FT_Error(0xB4):
		return ErrFontMissingCharsField
	case C.FT_Error(0xB5):
		return ErrFontMissingStartcharField
	case C.FT_Error(0xB6):
		return ErrFontMissingEncodingField
	case C.FT_Error(0xB7):
		return ErrFontMissingBbxField
	case C.FT_Error(0xB8):
		return ErrFontBbxTooBig
	case C.FT_Error(0xB9):
		return ErrFontCorruptedFontHeader
	case C.FT_Error(0xBA):
		return ErrFontCorruptedFontGlyphs
	default:
		return errors.New(fmt.Sprintf("Error Code 0x%04X", code))
	}
}
