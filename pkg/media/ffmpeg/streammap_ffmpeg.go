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

////////////////////////////////////////////////////////////////////////////////
// METHODS

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

// Map returns an output stream for an input stream
func (this *streammap) Map() map[*stream]*stream {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	result := make(map[*stream]*stream, len(this.m))
	for k, v := range this.m {
		result[k] = v
	}
	return result
}
