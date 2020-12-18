// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"net/url"
	"strconv"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ffmpeg "github.com/djthorpe/gopi/v3/pkg/sys/ffmpeg"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type outputctx struct {
	sync.RWMutex

	ctx     *ffmpeg.AVFormatContext
	avio    *ffmpeg.AVIOContext
	streams []*stream
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

	var result error

	// Write trailer
	if this.ctx != nil {
		if err := this.ctx.WriteTrailer(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Close files
	if this.avio != nil {
		this.avio.Flush()
		if err := this.avio.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release streams
	for _, stream := range this.streams {
		stream.Release()
	}

	// Close media
	if this.ctx != nil {
		this.ctx.Free()
	}

	// Release resources
	this.ctx = nil
	this.avio = nil
	this.streams = nil

	// Return success
	return multierror.Flatten(result)
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

func (this *outputctx) IsFile() bool {
	if this.ctx == nil {
		return false
	} else {
		return this.ctx.Flags()&ffmpeg.AVFMT_NOFILE == 0
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
	if this.IsFile() {
		flags |= gopi.MEDIA_FLAG_FILE
	}

	return flags
}

func (this *outputctx) Streams() []gopi.MediaStream {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Check for closed file
	if this.ctx == nil {
		return nil
	}

	// Return streams
	result := []gopi.MediaStream{}
	for _, stream := range this.streams {
		result = append(result, stream)
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

// DecodeIterator loops over selected streams from media object
func (this *outputctx) Write(ctx gopi.MediaDecodeContext, packet gopi.MediaPacket) error {
	// If file and no avio context, then create one
	if this.IsFile() && this.avio == nil {
		if avio, err := ffmpeg.NewAVIOContext(this.URL(), ffmpeg.AVIO_FLAG_WRITE); err != nil {
			return err
		} else {
			this.avio = avio
			this.ctx.SetIOContext(avio)
		}
	}

	// If streams have not been set up yet, create output streams and write
	// file header
	if this.streams == nil {
		if err := this.MapStreams(ctx); err != nil {
			return err
		}
		if err := this.ctx.WriteHeader(nil); err != nil {
			return err
		}
	}

	/*
		in_stream  = ifmt_ctx->streams[pkt.stream_index];
		if (pkt.stream_index >= stream_mapping_size ||
			stream_mapping[pkt.stream_index] < 0) {
			av_packet_unref(&pkt);
			continue;
		}
		pkt.stream_index = stream_mapping[pkt.stream_index];
		out_stream = ofmt_ctx->streams[pkt.stream_index];
		log_packet(ifmt_ctx, &pkt, "in");
		// copy packet
		pkt.pts = av_rescale_q_rnd(pkt.pts, in_stream->time_base, out_stream->time_base, AV_ROUND_NEAR_INF|AV_ROUND_PASS_MINMAX);
		pkt.dts = av_rescale_q_rnd(pkt.dts, in_stream->time_base, out_stream->time_base, AV_ROUND_NEAR_INF|AV_ROUND_PASS_MINMAX);
		pkt.duration = av_rescale_q(pkt.duration, in_stream->time_base, out_stream->time_base);
		pkt.pos = -1;
		log_packet(ofmt_ctx, &pkt, "out");
		ret = av_interleaved_write_frame(ofmt_ctx, &pkt);
		if (ret < 0) {
			fprintf(stderr, "Error muxing packet\n");
			break;
		}
		av_packet_unref(&pkt);
	*/

	return gopi.ErrNotImplemented
}

func (this *outputctx) MapStreams(ctx gopi.MediaDecodeContext) error {
	var mapper *streammap

	// Get stream map
	if ctx == nil {
		return gopi.ErrBadParameter.WithPrefix("Write")
	} else if ctx_, ok := ctx.(*decodectx); ok == false {
		return gopi.ErrBadParameter.WithPrefix("Write")
	} else {
		mapper = ctx_.streammap
	}

	// Iterate through streams and create output streams as necessary
	for in, out := range mapper.Map() {
		if out != nil {
			continue
		}
		// Create an output stream
		if stream := ffmpeg.NewStream(this.ctx, nil); stream == nil {
			return gopi.ErrBadParameter.WithPrefix("Write")
		} else if out = NewStream(stream, in); out == nil {
			return gopi.ErrInternalAppError.WithPrefix("Write")
		} else if err := mapper.Set(in, out); err != nil {
			return err
		} else {
			fmt.Printf("%v\n   => %v\n\n", in, out)
		}
	}

	// Success
	return nil
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
