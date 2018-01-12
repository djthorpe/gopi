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

type FontManager struct{}

type manager struct {
	log    gopi.Logger
	handle C.FT_Library
}

type ftError C.FT_Error

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FT_SUCCESS ftError = 0
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config FontManager) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.font.rpi.FontManager.Open>{ }")
	this := new(manager)
	this.log = log

	return this, nil
}

func (this *manager) Close() error {
	this.log.Debug("<sys.font.rpi.FontManager.Close>{ }")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *manager) String() string {
	return fmt.Sprintf("<sys.font.rpi.FontManager>{ }")
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE FUNCTIONS

func (this *manager) vgfontInit() ftError {
	if err := C.FT_Init_FreeType(&this.handle); ftError(err) != FT_SUCCESS {
		return ftError(err)
	} else {
		return FT_SUCCESS
	}
}

func (this *manager) vgfontDestroy() ftError {
	if err := C.FT_Done_FreeType(this.handle); ftError(err) != FT_SUCCESS {
		return ftError(err)
	} else {
		return FT_SUCCESS
	}
}
