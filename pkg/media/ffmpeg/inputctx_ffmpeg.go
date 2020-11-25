// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"io"
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
	ctx     *ffmpeg.AVFormatContext
	streams map[int]*stream
}

////////////////////////////////////////////////////////////////////////////////
// INIT AND CLOSE

func NewInputContext(ctx *ffmpeg.AVFormatContext) *inputctx {
	// Create object
	this := new(inputctx)
	if ctx == nil {
		return nil
	} else {
		this.ctx = ctx
	}

	// Create streams
	if streams := this.ctx.Streams(); streams == nil {
		return nil
	} else {
		this.streams = make(map[int]*stream, len(streams))
		for _, stream := range streams {
			key := stream.Index()
			this.streams[key] = NewStream(stream)
		}
	}

	// success
	return this
}

func (this *inputctx) Close() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Release resources
	this.streams = nil

	// Close media
	this.ctx.CloseInput()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *inputctx) URL() *url.URL {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.ctx.Url()
}

func (this *inputctx) Metadata() gopi.MediaMetadata {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return NewMetadata(this.ctx.Metadata())
}

func (this *inputctx) Flags() gopi.MediaFlag {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Stream flags
	flags := gopi.MEDIA_FLAG_FILE
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
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	result := []gopi.MediaStream{}
	for _, stream := range this.streams {
		result = append(result, stream)
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - ITERATE OVER PACKETS

// Iterate over packets in the input stream
func (this *inputctx) DecodeIterator(fn gopi.DecodeIteratorFunc) error {
	// Lock for writing as ReadPacket modifies state
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check parameters
	if fn == nil {
		return gopi.ErrBadParameter.WithPrefix("DecodeIterator")
	}

	// Create context
	ctx := NewDecodeContext()
	exists := false
	defer ctx.Close()

	for {
		if err := this.ctx.ReadPacket(ctx.packet); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if ctx.stream, exists = this.streams[ctx.packet.Stream()]; exists == false {
			return gopi.ErrInternalAppError
		} else if err := fn(ctx); err != nil {
			return err
		}
		ctx.Release()
	}

	// Return success
	return nil
}

/*
func (this *inputctx) FrameIterator(packet gopi.MediaPacket, fn FrameIteratorFunc) error {
	// We don't lock in this function, assuming locking is done via PacketIterator
	if fn == nil {
		return gopi.ErrBadParameter.WithPrefix("FrameIterator")
	}

	// Supply raw packet data as input to a decoder
	// https://ffmpeg.org/doxygen/trunk/group__lavc__decoding.html#ga58bc4bf1e0ac59e27362597e467efff3
	if err := avcodec_send_packet(pCodecContext, packet); err != nil {
		return err
	}

	for {
		if codec.Decode()
		frame.avcodec_receive_frame
		if err := this.ctx.ReadPacket(packet); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else if stream, exists := this.streams[packet.Stream()]; exists == false {
			return gopi.ErrInternalAppError
		} else if err := fn(packet, stream); err != nil {
			return err
		}
		packet.Release()
	}

	// Return success
	return nil
}
*/

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
