/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// VGFONT
//
// This package implements the ability to load TrueType fonts into OpenVG.
// interfacing with FreeType. It is based on the following sets of code:
//
//    https://github.com/ajstarks/openvg/blob/master/fontutil
//    https://github.com/raspberrypi/firmware/tree/master/opt/vc/src/hello_pi/libs/vgfont
//
// To construct a font library instance, use the following form:
//
//	fonts, err := gopi.Open(rpi.VGFont{ PPI: pixels_per_inch }, logger)
//	if err != nil { /* handle error */ }
//	defer fonts.Close()
//
// This returns a gopi.Driver object which can be later cast into a font driver
// of type khronos.VGFontDriver. You should pass the number of pixels per inch
// for your display into the configuration, which is used to calculate font
// size from points when drawing text. If you omit, then a default value will
// be set instead.
//
// To load a font face from a TrueType file, use the following:
//
//   face, err := fonts.(khronos.VGFontDriver).OpenFace(filename)
//   if err != nil { /* handle error */ }
//   defer fonts.DestroyFace(face)
//
// The face returned is of type kronos.VGFace. A single TrueType file may
// contain several faces. The number of faces can be returned as face.GetNumFaces()
// and you can load other faces in the same file using:
//
//   face, err := fonts.(khronos.VGFontDriver).OpenFaceAtIndex(filename,index)
//   if err != nil { /* handle error */ }
//   defer fonts.DestroyFace(face)
//
// You can enumerate the loaded font families and styles using the following:
//
//   for _, family := range fonts.(khronos.VGFontDriver).GetFamilies() {
//     styles := fonts.(khronos.VGFontDriver).GetFaces(family,khronos.VG_FONT_STYLE_ANY)
//     /* will return a map of style names to faces... */
//   }
//
// GetStyles can be called with an empty family string to return styles for all
// families. The following flags can be used to query family styles (the flags
// can be OR'd):
//
//   khronos.VG_FONT_STYLE_ANY : Any style
//   khronos.VG_FONT_STYLE_REGULAR : Regular style
//   khronos.VG_FONT_STYLE_BOLD : Bold style
//   khronos.VG_FONT_STYLE_ITALIC : Italic style
//
// To load many font faces at once, you can also import fonts from a folder. In
// order to do this, you'll need to define a callback function which can be used
// to determine if a particular font file should be loaded. For example:
//
//   func fontload_callback(filename string, info os.FileInfo) bool {
//     if strings.HasPrefix(info.Name(), ".") {
//       return false /* ignore hidden directories and files */
//     }
//     if info.IsDir() {
//       return true /* recurse into folders */
//     }
//     if path.Ext(filename) == ".ttf" || path.Ext(filename) == ".TTF" {
//       return true /* load files with TTF extension */
//     }
//     return false /* ignore all other files */
//   }
//
// Then you can load font faces in a folder using the following method:
//
//   err := fonts.(khronos.VGFontDriver).OpenFacesAtPath(basepath,fontload_callback)
//   if err != nil { /* handle errors */ }
//
// If there are several faces within a file, they will all be loaded.
//
// Further documentation on how to use the font faces for rendering will be
// provided shortly!
//
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
	util "github.com/djthorpe/gopi/util"
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
	PPI uint
}

// VGFont driver
type vgfDriver struct {
	handle C.FT_Library
	log    *util.LoggerDevice
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
}

// vgfEncoding represents charmap encoding
type vgfEncoding string

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_FONT_NONE                C.VGFont    = C.VGFont(0)
	VG_FONT_PPI_DEFAULT         uint        = 72
	VG_FONT_FT_ERROR_NONE       C.FT_Error  = C.FT_Error(0)
	VG_FONT_FT_ENCODING_UNICODE vgfEncoding = "unic"
	VG_FONT_FT_STYLE_FLAG_ITALIC C.FT_Long = (1 << 0)
	VG_FONT_FT_STYLE_FLAG_BOLD C.FT_Long = (1 << 1)
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

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

// Open the driver
func (config VGFont) Open(log *util.LoggerDevice) (gopi.Driver, error) {
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

	// Set PPI to default value if it's not set
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
	family_map := make(map[string]bool,0)
	family_array := make([]string,0)
	for _, face := range this.faces {
		family := face.GetFamily()
		if _, exists := family_map[family]; exists {
			continue
		}
		family_map[family] = true
		family_array = append(family_array,family)
	}
	return family_array
}

// Return faces in a family and/or with a particular set of attributes
func (this *vgfDriver) GetFaces(family string, flags khronos.VGFontStyleFlags) []khronos.VGFace {
	faces := make([]khronos.VGFace,0)
	for _, face := range this.faces {
		if family != "" && family != face.GetFamily() {
			continue
		}
		switch(flags) {
		case khronos.VG_FONT_STYLE_ANY:
			faces = append(faces,face)
			break
		case khronos.VG_FONT_STYLE_REGULAR:
			if face.IsBold() == false && face.IsItalic() == false {
				faces = append(faces,face)
			}
			break
		case khronos.VG_FONT_STYLE_BOLD:
			if face.IsBold() == true && face.IsItalic() == false {
				faces = append(faces,face)
			}
			break
		case khronos.VG_FONT_STYLE_ITALIC:
			if face.IsBold() == false && face.IsItalic() == true {
				faces = append(faces,face)
			}
			break
		case khronos.VG_FONT_STYLE_BOLDITALIC:
			if face.IsBold() == true && face.IsItalic() == true {
				faces = append(faces,face)
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
	return vgfontGetError(C.FT_Init_FreeType(unsafe.Pointer(&this.handle)))
}

func (this *vgfDriver) vgfontDestroy() error {
	return vgfontGetError(C.FT_Done_FreeType(unsafe.Pointer(this.handle)))
}

func (this *vgfDriver) vgfontLoadFace(path string, index uint) (C.FT_Face, error) {
	var face C.FT_Face
	ret := C.FT_New_Face(unsafe.Pointer(this.handle), C.CString(path), C.FT_Long(index), unsafe.Pointer(&face))
	return face, vgfontGetError(ret)
}

func (this *vgfDriver) vgfontDoneFace(handle C.FT_Face) error {
	return vgfontGetError(C.FT_Done_Face(unsafe.Pointer(handle)))
}

func (this *vgfDriver) vgfontSelectCharmap(handle C.FT_Face, encoding vgfEncoding) error {
	code := C.FT_Encoding(uint32(encoding[3]) | uint32(encoding[2])<<8 | uint32(encoding[1])<<16 | uint32(encoding[0])<<24)
	return vgfontGetError(C.FT_Select_Charmap(handle, code))
}

func (this *vgfDriver) vgfontLoadGlyphToFont(face *vgfFace, glyph C.FT_UInt) error {
	// create a path
	path := VG_PATH_NONE
	outline := face.handle.glyph.outline
	if outline.n_contours != 0 {
		path = C.vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_F, C.VGfloat(1.0), C.VGfloat(0.0), C.VGint(0), C.VGint(0), C.VGbitfield(VG_PATH_CAPABILITY_ALL))
		if path == VG_PATH_NONE {
			return vgGetError(vgErrorType(C.vgGetError()))
		}
	}

	/* TODO
	   	origin := []C.VGfloat{ C.VGfloat(0), C.VGfloat(0) }
	   	escapement :=
	         VGfloat escapement[] = { float_from_26_6(font->ft_face->glyph->advance.x), float_from_26_6(font->ft_face->glyph->advance.y) };
	         vgSetGlyphToPath(font->vg_font, glyph_index, vg_path, VG_FALSE, origin, escapement);
	*/

	// destroy path
	if path != VG_PATH_NONE {
		C.vgDestroyPath(path)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS: Errors

func vgfontGetError(code C.FT_Error) error {
	if code == C.FT_Error(0) {
		return nil
	}
	switch code {
	case C.FT_Error(0x01):
		return errors.New("Cannot Open Resource")
	case C.FT_Error(0x02):
		return errors.New("Unknown File Format")
	case C.FT_Error(0x03):
		return errors.New("Invalid File Format")
	case C.FT_Error(0x04):
		return errors.New("Invalid Freetype version")
	case C.FT_Error(0x05):
		return errors.New("Module Version is too low")
	case C.FT_Error(0x06):
		return errors.New("Invalid Argument")
	case C.FT_Error(0x07):
		return errors.New("Unimplemented Feature")
	default:
		return errors.New(fmt.Sprintf("Error Code 0x%04X", code))
	}
}
