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

	// Close media
	this.ctx.CloseInput()

	// Release resources
	for _, stream := range this.streams {
		stream.Release()
	}
	this.streams = nil
	this.ctx = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *inputctx) URL() *url.URL {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return nil
	} else {
		return this.ctx.Url()
	}
}

func (this *inputctx) Metadata() gopi.MediaMetadata {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.ctx == nil {
		return nil
	} else {
		return NewMetadata(this.ctx.Metadata())
	}
}

func (this *inputctx) Flags() gopi.MediaFlag {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Check for closed file
	if this.ctx == nil {
		return gopi.MEDIA_FLAG_NONE
	}

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

	// Check for closed file
	if this.ctx == nil {
		return nil
	}

	result := []gopi.MediaStream{}
	for _, stream := range this.streams {
		result = append(result, stream)
	}
	return result
}

func (this *inputctx) StreamsForFlag(flag gopi.MediaFlag) []int {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Check for closed file
	if this.ctx == nil {
		return nil
	}

	// Iterate through streams
	result := []int{}
	for key, stream := range this.streams {
		if stream.Flags()&flag != 0 {
			result = append(result, key)
		}
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - ITERATE OVER PACKETS

// Iterate over packets in the input stream
func (this *inputctx) DecodeIterator(streams []int, fn gopi.DecodeIteratorFunc) error {
	// Lock for writing as ReadPacket modifies state
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check parameters
	if fn == nil || this.ctx == nil {
		return gopi.ErrBadParameter.WithPrefix("DecodeIterator")
	}

	// If streams argument is empty or nil, select all streams
	if len(streams) == 0 {
		for index := range this.streams {
			streams = append(streams, index)
		}
	}

	// Create decode context map and call close on each on exit
	contextmap := make(map[int]*decodectx, len(this.streams))
	defer func() {
		for _, ctx := range contextmap {
			ctx.Close()
		}
	}()

	// Create decode contexts
	for _, index := range streams {
		if stream, exists := this.streams[index]; exists == false {
			return gopi.ErrInternalAppError.WithPrefix("DecodeIterator")
		} else if decodectx := NewDecodeContext(stream); decodectx == nil {
			return gopi.ErrInternalAppError.WithPrefix("DecodeIterator")
		} else {
			contextmap[index] = decodectx
		}
	}

	// Create a packet
	packet := ffmpeg.NewAVPacket()
	if packet == nil {
		return gopi.ErrInternalAppError.WithPrefix("DecodeIterator")
	}
	defer packet.Free()

	// Iterate over incoming packets, callback when packet should
	// be processed
	for {
		if err := this.ctx.ReadPacket(packet); err == io.EOF {
			// End of stream
			break
		} else if err != nil {
			return err
		} else if ctx, exists := contextmap[packet.Stream()]; exists {
			err := fn(ctx, packet)
			packet.Release()
			if err != nil {
				return err
			}
		}
	}

	// Return success
	return nil
}

func (this *inputctx) DecodeFrameIterator(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket, fn gopi.DecodeFrameIteratorFunc) error {
	// Check parameters
	if ctx == nil || fn == nil {
		return gopi.ErrBadParameter.WithPrefix("DecodeFrameIterator")
	}
	// Get internal context object and check more parameters
	ctx_, ok := ctx.(*decodectx)
	if ok == false || packet == nil {
		return gopi.ErrBadParameter.WithPrefix("DecodeFrameIterator")
	}

	// Lock context for writing
	//ctx_.RWMutex.Lock()
	//defer ctx_.RWMutex.Unlock()

	// Decode packet
	if err := ctx_.DecodePacket(packet); err != nil {
		return fmt.Errorf("DecodeFrameIterator: %w", err)
	}

	// Iterate through frames
	for {
		// Return frames until no more available
		if frame, err := ctx_.DecodeFrame(); err == io.EOF {
			return err
		} else if err != nil {
			return fmt.Errorf("DecodeFrameIterator: %w", err)
		} else if frame == nil {
			// Not enough data, so return without processing frame
			return nil
		} else {
			err := fn(ctx, frame)
			frame.Release()
			if err != nil {
				return err
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *inputctx) String() string {
	str := "<media"
	if url := this.URL(); url != nil {
		str += " url=" + strconv.Quote(url.String())
	}
	if metadata := this.Metadata(); metadata != nil {
		str += " metadata=" + fmt.Sprint(metadata)
	}
	if flags := this.Flags(); flags != gopi.MEDIA_FLAG_NONE {
		str += " flags=" + fmt.Sprint(flags)
	}
	if streams := this.Streams(); len(streams) > 0 {
		str += " streams=" + fmt.Sprint(streams)
	}
	return str + ">"
}
