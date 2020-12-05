// +build pulse

package pulse

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libpulse-simple
#include <pulse/simple.h>
#include <stdlib.h>
*/
import "C"

import (
	"time"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	PulseHandle     C.pa_simple
	PulseSampleSpec C.pa_sample_spec
	PulseChannelMap C.pa_channel_map
	PulseBufferAttr C.pa_buffer_attr
)

type (
	PulseStreamDirection int
	PulseSampleFormat    int
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	PA_STREAM_NODIRECTION PulseStreamDirection = iota
	PA_STREAM_PLAYBACK
	PA_STREAM_RECORD
	PA_STREAM_UPLOAD
)

const (
	PA_SAMPLE_U8 PulseSampleFormat = iota
	PA_SAMPLE_ALAW
	PA_SAMPLE_ULAW
	PA_SAMPLE_S16LE
	PA_SAMPLE_S16BE
	PA_SAMPLE_FLOAT32LE
	PA_SAMPLE_FLOAT32BE
	PA_SAMPLE_S32LE
	PA_SAMPLE_S32BE
	PA_SAMPLE_S24LE
	PA_SAMPLE_S24BE
	PA_SAMPLE_S24_32LE
	PA_SAMPLE_S24_32BE
	PA_SAMPLE_INVALID PulseSampleFormat = -1
)

////////////////////////////////////////////////////////////////////////////////
// PulseSampleSpec

func NewSampleSpec(fmt PulseSampleFormat, rate uint32, channels uint8) *PulseSampleSpec {
	this := new(PulseSampleSpec)
	ctx := (*C.pa_sample_spec)(this)
	ctx.format = (C.pa_sample_format_t)(fmt)
	ctx.rate = C.uint32_t(rate)
	ctx.channels = C.uint8_t(channels)
	return this
}

////////////////////////////////////////////////////////////////////////////////
// SIMPLE API

func PulseNewSimple(server, name string, dir PulseStreamDirection, dev, stream string, spec *PulseSampleSpec, channels *PulseChannelMap, attr *PulseBufferAttr) (*PulseHandle, error) {
	cServer := CString(server)
	cName := CString(name)
	cDev := CString(dev)
	cStream := CString(stream)
	defer C.free(unsafe.Pointer(cServer))
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cDev))
	defer C.free(unsafe.Pointer(cStream))

	var err C.int

	ctx := C.pa_simple_new(cServer, cName, C.pa_stream_direction_t(dir), cDev, cStream,
		(*C.pa_sample_spec)(spec),
		(*C.pa_channel_map)(channels),
		(*C.pa_buffer_attr)(attr),
		&err)
	if ctx == nil {
		return nil, PulseError(err)
	}

	// Return success
	return (*PulseHandle)(ctx), nil
}

func (this *PulseHandle) Free() {
	ctx := (*C.pa_simple)(this)
	C.pa_simple_free(ctx)
}

func (this *PulseHandle) Flush() error {
	var err C.int
	ctx := (*C.pa_simple)(this)
	if res := C.pa_simple_flush(ctx, &err); res != 0 {
		return PulseError(err)
	} else {
		return nil
	}
}

func (this *PulseHandle) Drain() error {
	var err C.int
	ctx := (*C.pa_simple)(this)
	if res := C.pa_simple_drain(ctx, &err); res != 0 {
		return PulseError(err)
	} else {
		return nil
	}
}

func (this *PulseHandle) GetLatency() (time.Duration, error) {
	var err C.int
	ctx := (*C.pa_simple)(this)
	if latency := C.pa_simple_get_latency(ctx, &err); err != 0 {
		return 0, PulseError(err)
	} else {
		return time.Duration(latency) * time.Microsecond, nil
	}
}

func (this *PulseHandle) Write(data []byte) error {
	return nil
}

func (this *PulseHandle) Read(data []byte) error {
	var err C.int
	ctx := (*C.pa_simple)(this)
	ptr := unsafe.Pointer(&data[0])
	size := len(data)
	if res := C.pa_simple_read(ctx, ptr, C.size_t(size), &err); res != 0 {
		return PulseError(err)
	} else {
		return nil
	}
}

func (this *PulseHandle) String() string {
	return "<pulse.simple>"
}

////////////////////////////////////////////////////////////////////////////////
// UTILS

func CString(value string) *C.char {
	if value == "" {
		return nil
	} else {
		return C.CString(value)
	}
}
