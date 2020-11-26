// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type decodectx struct {
	sync.RWMutex
	stream   *stream
	codecctx *ffmpeg.AVCodecContext
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

	// Return success
	return this
}

func (this *decodectx) Close() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Release resources
	this.stream = nil

	// Return success
	return nil
}

func (this *decodectx) Stream() gopi.MediaStream {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.stream
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
