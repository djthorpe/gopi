// +build freetype

package freetype

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ft "github.com/djthorpe/gopi/v3/pkg/sys/freetype"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type FontManager struct {
	gopi.Unit
	sync.Mutex

	library             ft.FT_Library
	major, minor, patch int
	faces               map[string]gopi.FontFace
}

////////////////////////////////////////////////////////////////////////////////
// INIT & CLOSE

func (this *FontManager) New(gopi.Config) error {
	if library, err := ft.FT_Init(); err != nil {
		return err
	} else {
		this.library = library
		this.major, this.minor, this.patch = ft.FT_Library_Version(library)
		this.faces = make(map[string]gopi.FontFace)
	}

	// Return success
	return nil
}

func (this *FontManager) Dispose() error {
	var result error

	for key, face := range this.faces {
		if err := this.DestroyFace(face); err != nil {
			result = multierror.Append(result, err)
		}
		delete(this.faces, key)
	}

	this.Lock()
	defer this.Unlock()

	if this.library != ft.FT_Library(nil) {
		if err := ft.FT_Destroy(this.library); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.library = nil
	this.faces = nil
	this.major, this.minor, this.patch = 0, 0, 0

	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *FontManager) String() string {
	str := "<fontmanager freetype"
	if this.major != 0 && this.minor != 0 && this.patch != 0 {
		str += fmt.Sprint(" version=", this.major, ".", this.minor, ".", this.patch)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION: gopi.FontManager

// Open a font face - first face at index 0 is loaded
func (this *FontManager) OpenFace(path string) (gopi.FontFace, error) {
	return this.OpenFaceAtIndex(path, 0)
}

// Open a font face - indexed within file of several faces
func (this *FontManager) OpenFaceAtIndex(path string, index uint) (gopi.FontFace, error) {
	this.Lock()
	defer this.Unlock()

	// Create the face
	face := NewFontFaceWithPath(path)

	if handle, err := ft.FT_NewFace(this.library, path, index); err != nil {
		return nil, err
	} else if err := ft.FT_SelectCharmap(handle, ft.FT_ENCODING_UNICODE); err != nil {
		ft.FT_DoneFace(handle)
		return nil, err
	} else {
		face.handle = handle
	}

	// Add face to list of faces
	this.faces[face.Path] = face

	return face, nil
}

// Open font faces at path, checking to see if individual files should
// be opened through a callback function
func (this *FontManager) OpenFacesAtPath(path string, callback func(manager gopi.FontManager, path string, info os.FileInfo) bool) error {
	if callback == nil {
		callback = openFacesAtPathDefaultCallback
	}
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if callback(this, path, info) == false {
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
		if face.NumFaces() > uint(1) {
			for i := uint(1); i < face.NumFaces(); i++ {
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

// Destroy a font face
func (this *FontManager) DestroyFace(face gopi.FontFace) error {
	this.Lock()
	defer this.Unlock()

	if face_, ok := face.(*fontface); ok == false || face_ == nil {
		return gopi.ErrBadParameter.WithPrefix("face")
	} else if _, exists := this.faces[face_.Path]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("face")
	} else {
		delete(this.faces, face_.Path)
		return ft.FT_DoneFace(face_.handle)
	}
}

// Return an array of font families which are loaded
func (this *FontManager) Families() []string {
	this.Lock()
	defer this.Unlock()

	families := make(map[string]bool, 0)
	for _, face := range this.faces {
		family := face.Family()
		if _, exists := families[family]; exists {
			continue
		}
		families[family] = true
	}
	familes_ := make([]string, 0, len(families))
	for k := range families {
		familes_ = append(familes_, k)
	}
	return familes_
}

// TODO: Return open face for filepath
func (this *FontManager) FaceForPath(path string) gopi.FontFace {
	return nil
}

// Return faces in a family and/or with a particular set of attributes
func (this *FontManager) Faces(family string, flags gopi.FontFlags) []gopi.FontFace {
	this.Lock()
	defer this.Unlock()

	faces := make([]gopi.FontFace, 0)
	for _, face := range this.faces {
		if family != "" && family != face.Family() {
			continue
		}
		switch flags {
		case gopi.FONT_FLAGS_STYLE_ANY:
			faces = append(faces, face)
		case gopi.FONT_FLAGS_STYLE_REGULAR, gopi.FONT_FLAGS_STYLE_BOLD, gopi.FONT_FLAGS_STYLE_ITALIC, gopi.FONT_FLAGS_STYLE_BOLDITALIC:
			if face.Flags()&flags == flags {
				faces = append(faces, face)
			}
		}
	}
	return faces
}

func openFacesAtPathDefaultCallback(_ gopi.FontManager, path string, info os.FileInfo) bool {

	// Ignore hidden files and folders
	if strings.HasPrefix(info.Name(), ".") {
		return false
	}

	// Allow recursing into any directory
	if info.IsDir() {
		return true
	}

	// Regular files: check supported file extensions
	if info.Mode().IsRegular() {
		if ext := strings.ToLower(filepath.Ext(info.Name())); ext == ".ttf" || ext == ".ttc" {
			return true
		} else if ext == ".otf" || ext == ".otc" {
			return true
		}
	}

	// Not supported
	return false
}
