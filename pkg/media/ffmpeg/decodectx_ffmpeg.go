// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"io"
	"sync"
	"syscall"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type decodectx struct {
	sync.RWMutex
	stream *stream
	ctx    *ffmpeg.AVCodecContext
	frame  *ffmpeg.AVFrame
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewDecodeContext(stream *stream) *decodectx {
	this := new(decodectx)

	// Check parameters
	if stream == nil {
		return nil
	} else {
		this.stream = stream
	}

	// Create frame
	if frame := ffmpeg.NewFrame(); frame == nil {
		return nil
	} else {
		this.frame = frame
	}

	// Create codec context
	if ctx := this.stream.NewContextWithOptions(nil); ctx == nil {
		return nil
	} else {
		this.ctx = ctx
	}

	// Return success
	return this
}

func (this *decodectx) Close() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Free context and frame
	this.frame.Free()
	this.ctx.Free()

	// Release resources
	this.stream = nil
	this.ctx = nil
	this.frame = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *decodectx) Stream() gopi.MediaStream {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.stream
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

// Present packet of data to the decoder
// https://ffmpeg.org/doxygen/trunk/group__lavc__decoding.html#ga58bc4bf1e0ac59e27362597e467efff3
func (this *decodectx) DecodePacket(packet gopi.MediaPacket) error {
	return this.ctx.DecodePacket(packet.(*ffmpeg.AVPacket))
}

// Decoded output data (into a frame) from a decoder
// https://ffmpeg.org/doxygen/trunk/group__lavc__decoding.html#ga11e6542c4e66d3028668788a1a74217c
func (this *decodectx) DecodeFrame() (*ffmpeg.AVFrame, error) {
	if err := this.ctx.DecodeFrame(this.frame); err == syscall.EAGAIN {
		// Not enough data
		return nil, nil
	} else if err == syscall.EINVAL {
		// End of stream
		return nil, io.EOF
	} else if err != nil {
		return nil, err
	} else {
		return this.frame, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *decodectx) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<mediacontext"
	if this.stream != nil {
		str += " stream=" + fmt.Sprint(this.stream)
	}
	return str + ">"
}
