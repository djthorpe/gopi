// +build ffmpeg

package ffmpeg

import (
	"sync"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type streammap struct {
	sync.RWMutex

	m map[*stream]*stream
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewStreamMap() *streammap {
	this := new(streammap)
	this.m = make(map[*stream]*stream)
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Add an input stream to the map
func (this *streammap) Set(in, out *stream) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if in == nil {
		return gopi.ErrBadParameter.WithPrefix("Set")
	} else if out == nil {
		this.m[in] = out
	} else if _, exists := this.m[in]; exists == false {
		return gopi.ErrNotFound.WithPrefix("Set")
	} else {
		this.m[in] = out
	}

	// Return success
	return nil
}

// Get returns an output stream for an input or
// nil otherwise
func (this *streammap) Get(in *stream) *stream {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if out, exists := this.m[in]; exists == false {
		return nil
	} else {
		return out
	}
}

// Map returns a map of output streams for a input
func (this *streammap) Map() map[*stream]*stream {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	result := make(map[*stream]*stream, len(this.m))
	for k, v := range this.m {
		result[k] = v
	}
	return result
}
