package ircodec

import (
	"context"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type IRCodec struct {
	gopi.Unit
	gopi.Logger
	gopi.Publisher
	sync.Mutex

	// List the codecs
	codecs []Codec
}

type Codec interface {
	Type() CodecType
	Process(gopi.LIRCEvent)
}

type CodecType uint

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	CODEC_NONE CodecType = iota
	CODEC_SONY_12
	CODEC_SONY_15
	CODEC_SONY_20
	CODEC_RC5_14
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *IRCodec) New(gopi.Config) error {
	this.codecs = append(this.codecs, NewRC5(CODEC_RC5_14))

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *IRCodec) Run(ctx context.Context) error {
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)
	for {
		select {
		case <-ctx.Done():
			return nil
		case evt := <-ch:
			if evt_, ok := evt.(gopi.LIRCEvent); ok {
				this.process(evt_)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROCESS EVENTS

func (this *IRCodec) process(evt gopi.LIRCEvent) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	this.Debug(evt.Mode(), " ", evt.Type(), " ", evt.Value(), "ms")
	for _, codec := range this.codecs {
		codec.Process(evt)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t CodecType) String() string {
	switch t {
	case CODEC_NONE:
		return "CODEC_NONE"
	case CODEC_SONY_12:
		return "CODEC_SONY_12"
	case CODEC_SONY_15:
		return "CODEC_SONY_15"
	case CODEC_SONY_20:
		return "CODEC_SONY_20"
	case CODEC_RC5_14:
		return "CODEC_RC5_14"
	default:
		return "[?? Invalid CodecType value]"
	}
}
