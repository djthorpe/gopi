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

type inputctx struct {
	sync.RWMutex
	ctx *ffmpeg.AVFormatContext
}

////////////////////////////////////////////////////////////////////////////////
// INIT AND CLOSE

func NewInputContext(ctx *ffmpeg.AVFormatContext) *inputctx {
	if ctx == nil {
		return nil
	} else {
		return &inputctx{ctx: ctx}
	}
}

func (this *inputctx) Close() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.ctx.CloseInput()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *inputctx) URL() *url.URL {
	return this.ctx.Url()
}

func (this *inputctx) Metadata() gopi.MediaMetadata {
	return NewMetadata(this.ctx.Metadata())
}

func (this *inputctx) Flags() gopi.MediaFlag {
	flags := gopi.MEDIA_FLAG_FILE

	// Stream flags
	for _, stream := range this.Streams() {
		flags |= stream.Flags()
	}

	// Add other flags with likely media file type
	metadata := this.Metadata()
	if flags&gopi.MEDIA_FLAG_AUDIO != 0 && metadata.Value(gopi.MEDIA_KEY_ALBUM) != nil {
		flags |= gopi.MEDIA_FLAG_ALBUM
	}
	if flags&gopi.MEDIA_FLAG_ALBUM != 0 && metadata.Value(gopi.MEDIA_KEY_ALBUM_ARTIST) != nil && metadata.Value(gopi.MEDIA_KEY_TITLE) != nil {
		flags |= gopi.MEDIA_FLAG_ALBUM_TRACK
	}
	if flags&gopi.MEDIA_FLAG_ALBUM != 0 {
		if compilation, ok := metadata.Value(gopi.MEDIA_KEY_COMPILATION).(bool); ok && compilation {
			flags |= gopi.MEDIA_FLAG_ALBUM_COMPILATION
		}
	}
	return flags
}

func (this *inputctx) Streams() []gopi.MediaStream {
	result := []gopi.MediaStream{}
	streams := this.ctx.Streams()
	if streams == nil {
		return nil
	}

	// Create stream array
	for _, stream := range streams {
		result = append(result, NewStream(stream))
	}

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *inputctx) String() string {
	str := "<media"
	if url := this.URL(); url != nil {
		str += " url=" + strconv.Quote(url.String())
	}
	str += " metadata=" + fmt.Sprint(this.Metadata())
	if flags := this.Flags(); flags != gopi.MEDIA_FLAG_NONE {
		str += " flags=" + fmt.Sprint(flags)
	}
	if streams := this.Streams(); len(streams) > 0 {
		str += " streams=" + fmt.Sprint(streams)
	}
	return str + ">"
}
