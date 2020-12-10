// +build chromaprint

package chromaprint

import (
	"fmt"
	"strconv"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	chromaprint "github.com/djthorpe/gopi/v3/pkg/sys/chromaprint"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type manager struct {
	gopi.Unit
	gopi.Logger
	sync.Mutex

	streams []*stream
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func (this *manager) New(gopi.Config) error {
	// Return success
	return nil
}

func (this *manager) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var result error
	for _, stream := range this.streams {
		if stream == nil {
			continue
		}
		if err := stream.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.streams = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *manager) String() string {
	str := "<chromaprint.manager"
	if v := chromaprint.Version(); v != "" {
		str += " version=" + strconv.Quote(v)
	}
	if len(this.streams) > 0 {
		str += " streams=" + fmt.Sprint(this.streams)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *manager) NewStream(rate, channels int) (Stream, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	stream, err := NewStream(rate, channels)
	if err != nil {
		return nil, err
	} else {
		this.streams = append(this.streams, stream)
	}

	// Return success
	return stream, nil
}

func (this *manager) Close(s Stream) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	s_, ok := s.(*stream)
	if s_ == nil || ok == false {
		return gopi.ErrBadParameter.WithPrefix("Close")
	}

	var result error
	for i := range this.streams {
		if s_ == this.streams[i] {
			if err := s_.Close(); err != nil {
				result = multierror.Append(result, err)
			}
			this.streams[i] = nil
		}
	}
	return result
}
