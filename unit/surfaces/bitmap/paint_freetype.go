// +build rpi,freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap

import (
	"image/color"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	ft "github.com/djthorpe/gopi/v2/sys/freetype"
)

////////////////////////////////////////////////////////////////////////////////
// RUNE

func (this *bitmap) PaintRunePx(c color.Color,pt gopi.Point,ch rune,face gopi.FontFace,pixels float32) {
	if image,err := face.BitmapForRunePixels(ch,pixels); err != nil {
		this.Log.Error(err)
	}
}

