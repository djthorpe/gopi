package dispmanx

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
#include <stdlib.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Data struct {
	buf uintptr
	cap uint32
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewData(size uint32) *Data {
	this := new(Data)
	this.cap = AlignUp(size, 4)
	this.buf = uintptr(C.malloc(C.uint(this.cap)))
	if this.cap == 0 || this.buf == 0 {
		return nil
	} else {
		return this
	}
}

func (this *Data) Dispose() {
	C.free(unsafe.Pointer(this.buf))
	this.buf = 0
	this.cap = 0
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

// Cap returns number of bytes capacity
func (this *Data) Cap() uint32 {
	return this.cap
}

// Ptr returns ptr to the data
func (this *Data) Ptr() uintptr {
	return this.buf
}

// PtrMinusOffset returns ptr to the data. minus a specific byte offset
func (this *Data) PtrMinusOffset(offset uint32) uintptr {
	return this.buf - uintptr(offset)
}

// Byte returns byte array from a specific byte offset
func (this *Data) Byte(offset uintptr) []byte {
	var result []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	hdr.Data = this.buf + uintptr(offset)
	hdr.Len = int(this.cap) - int(offset)
	hdr.Cap = hdr.Len
	return result
}

// Uint16 returns uint16 array from a specific byte offset
func (this *Data) Uint16(offset uintptr) []uint16 {
	var result []uint16
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	hdr.Data = this.buf + uintptr(offset)
	hdr.Len = (int(this.cap) - int(offset)) >> 1
	hdr.Cap = hdr.Len
	return result
}

// Uint32 returns uint32 array from a specific byte offset
func (this *Data) Uint32(offset uintptr) []uint32 {
	var result []uint32
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&result))
	hdr.Data = this.buf + uintptr(offset)
	hdr.Len = (int(this.cap) - int(offset)) >> 2
	hdr.Cap = hdr.Len
	return result
}

func (this *Data) SetUint8(offset uintptr, count uint32, data uint8) {
	buf := this.Byte(offset)
	for i := uint32(0); i < count; i++ {
		buf[i] = data
	}
}

func (this *Data) SetUint16(offset uintptr, count uint32, data uint16) {
	buf := this.Uint16(offset)
	for i := uint32(0); i < count; i++ {
		buf[i] = data
	}
}

func (this *Data) SetUint24(offset uintptr, count uint32, data uint32) {
	// We treat Uint24 as bytes
	buf := this.Byte(offset)
	ptr := 0
	for i := uint32(0); i < count; i++ {
		buf[ptr] = uint8(data)
		buf[ptr+1] = uint8(data >> 8)
		buf[ptr+2] = uint8(data >> 16)
		ptr += 3
	}
}

func (this *Data) SetUint32(offset uintptr, count uint32, data uint32) {
	buf := this.Uint32(offset)
	for i := uint32(0); i < count; i++ {
		buf[i] = data
	}
}

func (this *Data) GetUint8(offset uintptr) byte {
	ptr := (*byte)(unsafe.Pointer(this.buf + offset))
	return *ptr
}

func (this *Data) GetUint16(offset uintptr) uint16 {
	ptr := (*uint16)(unsafe.Pointer(this.buf + offset))
	return *ptr
}

func (this *Data) GetUint32(offset uintptr) uint32 {
	ptr := (*uint32)(unsafe.Pointer(this.buf + offset))
	return *ptr
}

func (this *Data) SetBit(offset uintptr, count uint32, bit bool) {
	buf := this.Byte(offset)
	bytecount := count >> 3
	for i := uint32(0); i < bytecount; i++ {
		if bit {
			buf[i] = 0xFF
		} else {
			buf[i] = 0x00
		}
	}
	if count&0x07 != 0 {
		bitmask := uint8(1<<(8-count&0x07) - 1) // Bitmask is which bits to carry over
		if bit {
			buf[bytecount] = buf[bytecount] | ^bitmask // Set bits
		} else {
			buf[bytecount] = buf[bytecount] & bitmask // Clear bits
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Data) String() string {
	str := "<data"
	str += fmt.Sprint(" size=", this.cap)
	str += fmt.Sprint(" data=", strings.ToUpper(hex.EncodeToString(this.Byte(0))))
	return str + ">"
}
