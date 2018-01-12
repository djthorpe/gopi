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

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
  #cgo CFLAGS:   -I/usr/include/freetype2
  #cgo LDFLAGS:  -lfreetype
  #include <ft2build.h>
  #include <freetype/freetype.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type FontManager struct {
	// RootPath for relative OpenFace calls
	RootPath string
}

type manager struct {
	log       gopi.Logger
	root_path string
	lock      sync.Mutex
	library   ftLibrary
}

type (
	ftError   C.FT_Error
	ftLibrary C.FT_Library
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FT_SUCCESS ftError = 0
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
	return this, nil
}

func (this *manager) Close() error {
	this.log.Debug("<sys.font.rpi.FontManager.Close>{ }")

	this.lock.Lock()
	defer this.lock.Unlock()

	if this.library == FT_NO_LIBRARY {
		return nil
	} else if err := ftDestroy(this.library); err != FT_SUCCESS {
		this.library = nil
		return err
	} else {
		this.library = nil
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *manager) String() string {
	return fmt.Sprintf("<sys.font.rpi.FontManager>{ }")
}

func (e ftError) Error() string {
	switch e {
	case FT_SUCCESS:
		return "FT_SUCCESS"
	default: // TODO Widen out error messages!
		return "FT_ERROR"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS

func ftInit() (ftLibrary, ftError) {
	var handle C.FT_Library
	if err := C.FT_Init_FreeType(&handle); ftError(err) != FT_SUCCESS {
		return FT_NO_LIBRARY, ftError(err)
	} else {
		return ftLibrary(handle), FT_SUCCESS
	}
}

func ftDestroy(handle ftLibrary) ftError {
	if err := C.FT_Done_FreeType(handle); ftError(err) != FT_SUCCESS {
		return ftError(err)
	} else {
		return FT_SUCCESS
	}
}
