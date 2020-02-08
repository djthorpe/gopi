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
)

////////////////////////////////////////////////////////////////////////////////
// RUNE

func (this *bitmap) Rune(c color.Color, origin gopi.Point, ch rune, face gopi.FontFace, size gopi.FontSize) {
	if image, err := face.BitmapForRunePixels(ch, 32); err != nil {
		this.Log.Error(err)
	} else {
		fmt.Println(image)
	}
}
