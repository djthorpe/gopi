// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"image"
	"image/color"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type frame struct {
	ctx    *ffmpeg.AVFrame
	planes [][]uint8
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewFrame() *frame {
	if ctx := ffmpeg.NewFrame(); ctx == nil {
		return nil
	} else {
		return &frame{ctx, nil}
	}
}

func (this *frame) Retain() error {
	return gopi.ErrNotImplemented
}

func (this *frame) Release() {
	this.planes = nil
	this.ctx.Release()
}

func (this *frame) Free() {
	this.Release()
	if this.ctx != nil {
		this.ctx.Free()
		this.ctx = nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMAGE IMPLEMENTATION

func (this *frame) ColorModel() color.Model {
	return color.YCbCrModel
}

func (this *frame) Bounds() image.Rectangle {
	return image.Rectangle{
		image.ZP,
		image.Pt(this.ctx.PictSize()),
	}
}

func (this *frame) Bytes(plane int) []byte {
	if this.planes == nil || plane < 0 || plane >= len(this.planes) {
		return nil
	}
	return this.planes[plane]
}

func (this *frame) Stride(plane int) int {
	if this.planes == nil || plane < 0 || plane >= len(this.planes) {
		return 0
	}
	return this.ctx.StrideForPlane(plane)
}

func (this *frame) At(x, y int) color.Color {
	strideY := this.Stride(0)
	strideCb := this.Stride(1)
	strideCr := this.Stride(2)
	Y := this.Bytes(0)[x+y*strideY]
	Cb := this.Bytes(1)[x>>1+y*strideCb]
	Cr := this.Bytes(2)[x>>1+y*strideCr]
	return color.YCbCr{Y, Cb, Cr}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *frame) String() string {
	if this.ctx == nil {
		return "nil"
	} else {
		return fmt.Sprint(this.ctx)
	}
}
