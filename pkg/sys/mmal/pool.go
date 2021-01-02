//+build mmal

package mmal

import (
	"fmt"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>

// Callback Functions
MMAL_BOOL_T mmal_pool_callback(MMAL_POOL_T* pool, MMAL_BUFFER_HEADER_T* buffer,void* userdata);

static void mmal_pool_callback_set_ex(MMAL_POOL_T* pool,void* userdata) {
	mmal_pool_callback_set(pool,mmal_pool_callback,userdata);
}
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALPool  (C.MMAL_POOL_T)
	MMALQueue (C.MMAL_QUEUE_T)
)

type MMALPoolCallback func(pool *MMALPool, buffer *MMALBuffer, userdata uintptr) bool

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *MMALPort) CreatePool(num, size uint32) *MMALPool {
	ctx := (*C.MMAL_PORT_T)(this)
	if pool := C.mmal_port_pool_create(ctx, C.uint32_t(num), C.uint32_t(size)); pool == nil {
		return nil
	} else {
		ctx.buffer_num = C.uint32_t(num)
		ctx.buffer_size = C.uint32_t(size)
		C.mmal_pool_callback_set_ex(pool, nil)
		return (*MMALPool)(pool)
	}
}

func (this *MMALPort) FreePool(pool *MMALPool) {
	ctx := (*C.MMAL_PORT_T)(this)
	C.mmal_port_pool_destroy(ctx, (*C.MMAL_POOL_T)(pool))
}

func (this *MMALPool) SetCallback(fn MMALPoolCallback, userdata uintptr) {
	ctx := (*C.MMAL_POOL_T)(this)
	if fn != nil {
		C.mmal_pool_callback_set_ex(ctx, unsafe.Pointer(userdata))
		MMALPoolRegisterCallback(ctx, fn)
	} else {
		MMALPoolRegisterCallback(ctx, nil)
	}
}

func MMALQueueCreate() *MMALQueue {
	return (*MMALQueue)(C.mmal_queue_create())
}

func (this *MMALQueue) Free() {
	ctx := (*C.MMAL_QUEUE_T)(this)
	C.mmal_queue_destroy(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// CALLBACK REGISTRATION

var (
	poolCallback = make(map[*C.MMAL_POOL_T]MMALPoolCallback)
)

func MMALPoolRegisterCallback(pool *C.MMAL_POOL_T, fn MMALPoolCallback) {
	if fn != nil {
		poolCallback[pool] = fn
	} else {
		delete(poolCallback, pool)
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS - POOL

func (this *MMALPool) Get() *MMALBuffer {
	return this.Q().Get()
}

func (this *MMALPool) Put(buffer *MMALBuffer) {
	this.Q().Put(buffer)
}

func (this *MMALPool) PutBack(buffer *MMALBuffer) {
	this.Q().PutBack(buffer)
}

func (this *MMALPool) Count() uint32 {
	ctx := (*C.MMAL_POOL_T)(this)
	return uint32(ctx.headers_num)
}

func (this *MMALPool) Q() *MMALQueue {
	ctx := (*C.MMAL_POOL_T)(this)
	return (*MMALQueue)(ctx.queue)
}

func (this *MMALPool) Resize(num, size uint32) error {
	ctx := (*C.MMAL_POOL_T)(this)
	if err := Error(C.mmal_pool_resize(ctx, C.uint32_t(num), C.uint32_t(size))); err == MMAL_SUCCESS {
		return nil
	} else {
		return err
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS - QUEUE

func (this *MMALQueue) Length() uint {
	ctx := (*C.MMAL_QUEUE_T)(this)
	return uint(C.mmal_queue_length(ctx))
}

func (this *MMALQueue) Get() *MMALBuffer {
	ctx := (*C.MMAL_QUEUE_T)(this)
	return (*MMALBuffer)(C.mmal_queue_get(ctx))
}

func (this *MMALQueue) Put(buffer *MMALBuffer) {
	ctx := (*C.MMAL_QUEUE_T)(this)
	C.mmal_queue_put(ctx, (*C.MMAL_BUFFER_HEADER_T)(buffer))
}

func (this *MMALQueue) PutBack(buffer *MMALBuffer) {
	ctx := (*C.MMAL_QUEUE_T)(this)
	C.mmal_queue_put_back(ctx, (*C.MMAL_BUFFER_HEADER_T)(buffer))
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *MMALPool) String() string {
	str := "<mmal.pool"
	if count := this.Count(); count != 0 {
		str += " count=" + fmt.Sprint(count)
	}
	if q := this.Q(); q != nil {
		str += " q=" + fmt.Sprint(q)
	}
	return str + ">"
}

func (this *MMALQueue) String() string {
	str := "<mmal.queue"
	if length := this.Length(); length != 0 {
		str += " length=" + fmt.Sprint(length)
	}
	return str + ">"
}

//export mmal_pool_callback
func mmal_pool_callback(ctx *C.MMAL_POOL_T, bufferctx *C.MMAL_BUFFER_HEADER_T, userdata unsafe.Pointer) C.MMAL_BOOL_T {
	pool := (*MMALPool)(ctx)
	buffer := (*MMALBuffer)(bufferctx)
	if cb, exists := poolCallback[ctx]; exists {
		return mmal_from_bool(cb(pool, buffer, uintptr(userdata)))
	} else {
		// Empty buffer is available - queue it
		pool.Put(buffer)
		fmt.Println("mmal_pool_callback: Queued empty buffer:",pool)
		return MMAL_BOOL_FALSE
	}
}
