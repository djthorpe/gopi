// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: freetype2
#include <ft2build.h>
#include FT_FREETYPE_H
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// FACE FUNCTIONS

func FT_BitmapSize(handle FT_Bitmap) (uint, uint) {
	return uint(handle.width), uint(handle.rows)
}

func FT_BitmapPitch(handle FT_Bitmap) int {
	return int(handle.pitch)
}


func FT_BitmapStride(handle FT_Bitmap) uint {
	if int(handle.pitch) < 0 {
		return uint(-handle.pitch)
	} else {
		return uint(handle.pitch)
	}
}

func FT_BitmapNumGrays(handle FT_Bitmap) uint {
	if FT_PixelMode(handle.pixel_mode) == FT_PIXEL_MODE_GRAY {
		return uint(handle.num_grays)
	} else {
		return 0
	}
}

func FT_BitmapPixelMode(handle FT_Bitmap) FT_PixelMode {
	return FT_PixelMode(handle.pixel_mode)
}

func FT_BitmapBufferSize(handle FT_Bitmap) uint {
	_, h := FT_BitmapSize(handle)
	return FT_BitmapStride(handle) * h
}

func FT_BitmapData(handle FT_Bitmap) []byte {
	var data []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	hdr.Data = uintptr(unsafe.Pointer(handle.buffer))
	hdr.Len = int(FT_BitmapBufferSize(handle))
	hdr.Cap = hdr.Len
	return data
}

// FT_BitmapPixelsForRow returns pixel value array with elements between 0 and maximum for bitmap
// (depending on pixel mode) or nil if the y value is out of range or some other error
func FT_BitmapPixelsForRow(handle FT_Bitmap,y uint) []uint32 {
	if handle.buffer == nil {
		return nil
	}
	if y >= uint(handle.rows) {
		return nil
	}
	pixels := make([]uint32,int(handle.width))

	// Calculate row offset
	ptr := uintptr(unsafe.Pointer(handle.buffer)) + uintptr(y * FT_BitmapStride(handle))

	// Bits per pixel and mask
	bits_per_pixel := FT_BitmapBitsPerPixel(handle)
	mask := byte(1 << bits_per_pixel) - 1

	// Iterate along row
	shift := uint(0)
	mode := FT_PixelMode(handle.pixel_mode)
	for x := uint(0); x < uint(handle.width); x++ {
		if mode == FT_PIXEL_MODE_BGRA {
			data := (*uint32)(unsafe.Pointer(ptr))
			pixels[x] = (*data) // Elements in B,G,R,A order
		} else {
			offset := uintptr(bits_per_pixel * x >> 3)
			data := (*byte)(unsafe.Pointer(ptr+offset))
			if shift == 0 {
				pixels[x] = uint32(*data & mask)
			} else {
				pixels[x] = uint32((*data >> (7 - shift)) & mask)
			}
			// TODO fmt.Println("mode",FT_PixelMode(handle.pixel_mode),"offset",offset,"shift",shift,"pixel[",x,"]=",pixels[x],"mask",mask)
			shift = (shift + bits_per_pixel) % 8
		}
	}
	return pixels
}

func FT_BitmapBitsPerPixel(handle FT_Bitmap) uint {
	switch FT_PixelMode(handle.pixel_mode) {
	case FT_PIXEL_MODE_MONO:
		return 1 // 8 pixels per byte, mask = 0x01
	case FT_PIXEL_MODE_GRAY:
		return 8 // 1 pixel per byte, mask = 0xFF
	case FT_PIXEL_MODE_GRAY2:
		return 2 // 4 pixels per byte, mask = 0x03
	case FT_PIXEL_MODE_GRAY4:
		return 4 // 2 pixels per byte, mask = 0x0F
	case FT_PIXEL_MODE_LCD, FT_PIXEL_MODE_LCD_V:
		return 8 // 1 pixel per byte, mask = 0xFF
	case FT_PIXEL_MODE_BGRA:
		return 32 // 4 bytes per pixel, no mask
	default:
		return 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (handle FT_Bitmap) String() string {
	w, h := FT_BitmapSize(handle)
	str := "<FT_Bitmap "
	str += "size={" + fmt.Sprint(w) + "," + fmt.Sprint(h) + "} "
	str += "pixel_mode=" + fmt.Sprint(FT_BitmapPixelMode(handle)) + " "
	if FT_BitmapPixelMode(handle) == FT_PIXEL_MODE_GRAY {
		str += "num_grays=" + fmt.Sprint(FT_BitmapNumGrays(handle)) + " "
	}
	str += "pitch=" + fmt.Sprint(FT_BitmapPitch(handle)) + " "
	str += "data=" + strings.ToUpper(hex.EncodeToString(FT_BitmapData(handle)))
	return str + ">"
}
