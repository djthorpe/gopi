// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type outputctx struct {
	sync.RWMutex
	ctx *ffmpeg.AVFormatContext
}

////////////////////////////////////////////////////////////////////////////////
// INIT AND CLOSE

func NewOutputContext(ctx *ffmpeg.AVFormatContext) *outputctx {
	// Create object
	this := new(outputctx)
	if ctx == nil {
		return nil
	} else {
		this.ctx = ctx
	}

	// success
	return this
}

func (this *outputctx) Close() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Close media
	if this.ctx != nil {
		this.ctx.Free()
	}

	// Release resources
	this.ctx = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *outputctx) URL() *url.URL {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return nil
	} else {
		return this.ctx.Url()
	}
}

func (this *outputctx) Metadata() gopi.MediaMetadata {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return nil
	} else {
		return NewMetadata(this.ctx.Metadata())
	}
}

func (this *outputctx) Flags() gopi.MediaFlag {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Check for closed file
	if this.ctx == nil {
		return gopi.MEDIA_FLAG_NONE
	}

	// Stream flags
	flags := gopi.MEDIA_FLAG_ENCODER
	if this.ctx.Flags()&ffmpeg.AVFMT_NOFILE == 0 {
		flags |= gopi.MEDIA_FLAG_FILE
	}

	return flags
}

func (this *outputctx) Streams() []gopi.MediaStream {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

// DecodeIterator loops over selected streams from media object
func (this *outputctx) Write(gopi.MediaDecodeContext, gopi.MediaPacket) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *outputctx) String() string {
	str := "<ffmpeg.media"
	if url := this.URL(); url != nil {
		str += " output_url=" + strconv.Quote(url.String())
	}
	if metadata := this.Metadata(); metadata != nil {
		str += " metadata=" + fmt.Sprint(metadata)
	}
	if flags := this.Flags(); flags != 0 {
		str += " flags=" + fmt.Sprint(flags)
	}
	if streams := this.Streams(); len(streams) > 0 {
		str += " streams=" + fmt.Sprint(streams)
	}
	return str + ">"
}
