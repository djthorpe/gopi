// +build pulse

package pulse

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libpulse
#include <pulse/thread-mainloop.h>
#include <stdlib.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	PulseMainloop C.pa_threaded_mainloop
)

////////////////////////////////////////////////////////////////////////////////
// MAINLOOP METHODS

func NewMainloop() *PulseMainloop {
	if ctx := C.pa_threaded_mainloop_new(); ctx == nil {
		return nil
	} else {
		return (*PulseMainloop)(ctx)
	}
}

func (this *PulseMainloop) Free() {
	ctx := (*C.pa_threaded_mainloop)(this)
	C.pa_threaded_mainloop_free(ctx)
}

func (this *PulseMainloop) Start() error {
	ctx := (*C.pa_threaded_mainloop)(this)
	if ret := C.pa_threaded_mainloop_start(ctx); ret != 0 {
		return PulseError(ret)
	} else {
		return nil
	}
}

func (this *PulseMainloop) Stop() {
	ctx := (*C.pa_threaded_mainloop)(this)
	C.pa_threaded_mainloop_stop(ctx)
}

func (this *PulseMainloop) Lock() {
	ctx := (*C.pa_threaded_mainloop)(this)
	C.pa_threaded_mainloop_lock(ctx)
}

func (this *PulseMainloop) Unlock() {
	ctx := (*C.pa_threaded_mainloop)(this)
	C.pa_threaded_mainloop_unlock(ctx)
}
