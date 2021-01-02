//+build mmal

package mmal

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_connection.h>
*/
import "C"
import (
	"fmt"
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALConnection      C.MMAL_CONNECTION_T
	MMALConnectionFlags (C.uint32_t)
)

const (
	MMAL_CONNECTION_FLAG_TUNNELLING               MMALConnectionFlags = C.MMAL_CONNECTION_FLAG_TUNNELLING
	MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT      MMALConnectionFlags = C.MMAL_CONNECTION_FLAG_ALLOCATION_ON_INPUT
	MMAL_CONNECTION_FLAG_ALLOCATION_ON_OUTPUT     MMALConnectionFlags = C.MMAL_CONNECTION_FLAG_ALLOCATION_ON_OUTPUT
	MMAL_CONNECTION_FLAG_KEEP_BUFFER_REQUIREMENTS MMALConnectionFlags = C.MMAL_CONNECTION_FLAG_KEEP_BUFFER_REQUIREMENTS
	MMAL_CONNECTION_FLAG_DIRECT                   MMALConnectionFlags = C.MMAL_CONNECTION_FLAG_DIRECT
	MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS        MMALConnectionFlags = C.MMAL_CONNECTION_FLAG_KEEP_PORT_FORMATS
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func MMALConnectionCreate(out, in *MMALPort, flags MMALConnectionFlags) (*MMALConnection, error) {
	var ctx (*C.MMAL_CONNECTION_T)
	if err := Error(C.mmal_connection_create(&ctx, (*C.MMAL_PORT_T)(out), (*C.MMAL_PORT_T)(in), C.uint32_t(flags))); err == MMAL_SUCCESS {
		return (*MMALConnection)(ctx), nil
	} else {
		return nil, err
	}
}

func (this *MMALConnection) Free() error {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	if err := Error(C.mmal_connection_destroy(ctx)); err == MMAL_SUCCESS {
		return nil
	} else {
		return err
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *MMALConnection) Acquire() {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	C.mmal_connection_acquire(ctx)
}

func (this *MMALConnection) Release() error {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	if err := Error(C.mmal_connection_release(ctx)); err == MMAL_SUCCESS {
		return nil
	} else {
		return err
	}
}

func (this *MMALConnection) Enable() error {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	if err := Error(C.mmal_connection_enable(ctx)); err == MMAL_SUCCESS {
		return nil
	} else {
		return err
	}
}

func (this *MMALConnection) Disable() error {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	if err := Error(C.mmal_connection_disable(ctx)); err == MMAL_SUCCESS {
		return nil
	} else {
		return err
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *MMALConnection) Enabled() bool {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	return ctx.is_enabled != 0
}

func (this *MMALConnection) Name() string {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	return C.GoString(ctx.name)
}

func (this *MMALConnection) In() *MMALPort {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	return (*MMALPort)(ctx.in)
}

func (this *MMALConnection) Out() *MMALPort {
	ctx := (*C.MMAL_CONNECTION_T)(this)
	return (*MMALPort)(ctx.out)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *MMALConnection) String() string {
	str := "<mmal.connection"
	str += " enabled=" + fmt.Sprint(this.Enabled())
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if in := this.In(); in != nil {
		str += " in=" + strconv.Quote(in.Name())
	}
	if out := this.Out(); out != nil {
		str += " out=" + strconv.Quote(out.Name())
	}
	return str + ">"
}
