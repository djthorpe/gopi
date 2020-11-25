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
	packet *ffmpeg.AVPacket
	stream *stream
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewDecodeContext() *decodectx {
	this := new(decodectx)

	if packet := ffmpeg.NewAVPacket(); packet == nil {
		return nil
	} else {
		this.packet = packet
	}

	// Return success
	return this
}

func (this *decodectx) Close() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Free packet
	this.packet.Free()

	// Release resources
	this.packet = nil
	this.stream = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *decodectx) Bytes() []byte {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.packet == nil {
		return nil
	} else {
		return this.packet.Bytes()
	}
}

func (this *decodectx) Stream() gopi.MediaStream {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.stream
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *decodectx) Release() {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.packet != nil {
		this.packet.Release()
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *decodectx) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<mediacontext"
	if this.packet != nil {
		str += " packet=" + fmt.Sprint(this.packet)
	}
	if this.stream != nil {
		str += " stream=" + fmt.Sprint(this.stream)
	}
	return str + ">"
}
