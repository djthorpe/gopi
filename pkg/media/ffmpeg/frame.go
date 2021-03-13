// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"image"
	"image/color"

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
	if ctx := ffmpeg.NewAVFrame(); ctx == nil {
		return nil
	} else {
		return &frame{ctx, nil}
	}
}

func (this *frame) Retain() error {
	// To retain the frame, create the arrays of planes of data
	this.planes = nil
	i := 0
	for {
		if buf := this.ctx.Buffer(i); buf == nil {
			break
		} else {
			this.planes = append(this.planes, buf.Data())
		}
		i++
	}

	// Return success
	return nil
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
		image.Pt(this.ctx.PictWidth(), this.ctx.PictHeight()),
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
	// Currently assumes YUV420P
	Y := this.Bytes(0)[x+y*strideY]
	Cb := this.Bytes(1)[x>>1+y>>1*strideCb]
	Cr := this.Bytes(2)[x>>1+y>>1*strideCr]
	return color.YCbCr{Y, Cb, Cr}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *frame) String() string {
	str := "<MediaFrame"
	if this.ctx != nil {
		str += fmt.Sprint(" type=", this.ctx)
		str += fmt.Sprint(" bounds=", this.Bounds())
	}
	return str + ">"
}
