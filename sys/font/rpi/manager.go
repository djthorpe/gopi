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
// TYPES

type FontManager struct{}

type manager struct {
	log gopi.Logger
	handle C.FT_Library
}

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
  #cgo CFLAGS:   -I/usr/include/freetype2 -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lfreetype
  #include <ft2build.h>
  #include <freetype.h>
*/
import "C"

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

func (this *manager) vgfontInit() error {
	return vgfontGetError(C.FT_Init_FreeType(&this.handle))
}

func (this *manager) vgfontDestroy() error {
	return vgfontGetError(C.FT_Done_FreeType(this.handle))
}
