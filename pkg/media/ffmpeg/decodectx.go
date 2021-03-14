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

	stream    *stream
	frame     *frame
	ctx       *ffmpeg.AVCodecContext
	streammap *streammap
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewDecodeContext(s *stream, m *streammap) *decodectx {
	this := new(decodectx)

	// Check parameters
	if s == nil || m == nil {
		return nil
	} else {
		this.stream = s
		this.streammap = m
	}

	// Create frame
	if frame := NewFrame(); frame == nil {
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
	this.streammap = nil
	this.ctx = nil
	this.frame = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

// Stream returns stream associated with decode or nil
func (this *decodectx) Stream() gopi.MediaStream {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return nil
	} else {
		return this.stream
	}
}

// Frame returns current frame number or -1
func (this *decodectx) Frame() int {
	if this.ctx == nil {
		return -1
	} else {
		return this.ctx.Frame()
	}
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
func (this *decodectx) DecodeFrame() (*frame, error) {
	if err := this.ctx.DecodeFrame(this.frame.ctx); err == syscall.EAGAIN {
		// Not enough data
		return nil, nil
	} else if err == syscall.EINVAL {
		// End of stream
		return nil, io.EOF
	} else if err != nil {
		return nil, err
	} else if err := this.frame.Retain(); err != nil {
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

	str := "<ffmpeg.mediacontext"
	if frame_number := this.Frame(); frame_number >= 0 {
		str += " frame_number=" + fmt.Sprint(frame_number)
	}
	if this.ctx != nil {
		str += " " + fmt.Sprint(this.ctx)
	}
	return str + ">"
}
