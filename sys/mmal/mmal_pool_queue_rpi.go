//+build mmal

package mmal

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>

// Callback Functions
MMAL_BOOL_T mmal_pool_callback(MMAL_POOL_T* pool, MMAL_BUFFER_HEADER_T* buffer,void* userdata);
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CALLBACK REGISTRATION

var (
	pool_callback = make(map[*C.MMAL_POOL_T]MMAL_PoolCallback, 0)
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - POOLS

func MMALPortPoolCreate(handle MMAL_PortHandle, num, payload_size uint32, callback MMAL_PoolCallback, userdata uintptr) (MMAL_Pool, error) {
	if pool := C.mmal_port_pool_create(handle, C.uint32_t(num), C.uint32_t(payload_size)); pool == nil {
		return nil, MMAL_EINVAL
	} else {
		C.mmal_pool_callback_set(pool, C.MMAL_POOL_BH_CB_T(C.mmal_pool_callback), unsafe.Pointer(userdata))
		if callback != nil {
			pool_callback[pool] = callback
		}
		return pool, nil
	}
}

func MMALPortPoolDestroy(handle MMAL_PortHandle, pool MMAL_Pool) {
	if _, exists := pool_callback[pool]; exists {
		delete(pool_callback, pool)
	}
	C.mmal_port_pool_destroy(handle, pool)
}

func MMALPoolGetBuffer(pool MMAL_Pool) MMAL_Buffer {
	return MMALQueueGet(pool.queue)
}

func MMALPoolPutBuffer(pool MMAL_Pool, buffer MMAL_Buffer) {
	MMALQueuePut(pool.queue, buffer)
}

func MMALPoolReleaseBuffer(buffer MMAL_Buffer) {
	MMALBufferRelease(buffer)
}

func MMALPoolResize(handle MMAL_Pool, num, payload_size uint32) error {
	if status := MMAL_Status(C.mmal_pool_resize(handle, C.uint32_t(num), C.uint32_t(payload_size))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPoolString(pool MMAL_Pool) string {
	if pool == nil {
		return "<MMAL_Pool>{ nil }"
	} else {
		buffers := mmal_pool_buffer_array(pool)
		buffers_string := ""
		for _, buffer := range buffers {
			buffers_string += MMALBufferString(buffer) + " "
		}
		return fmt.Sprintf("<MMAL_Pool>{ queue=%v buffers=[ %v] }", MMALQueueString(pool.queue), buffers_string)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - QUEUES

func MMALQueueCreate() MMAL_Queue {
	return MMAL_Queue(C.mmal_queue_create())
}

func MMALQueueDestroy(handle MMAL_Queue) {
	C.mmal_queue_destroy(handle)
}

func MMALQueueString(handle MMAL_Queue) string {
	if handle == nil {
		return "<MMAL_Queue>{ nil }"
	} else {
		return fmt.Sprintf("<MMAL_Queue>{ length=%v }", C.mmal_queue_length(handle))
	}
}

func MMALQueuePut(handle MMAL_Queue, buffer MMAL_Buffer) {
	C.mmal_queue_put(handle, buffer)
}

func MMALQueueGet(handle MMAL_Queue) MMAL_Buffer {
	return MMAL_Buffer(C.mmal_queue_get(handle))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func mmal_pool_buffer_array(pool MMAL_Pool) []MMAL_Buffer {
	var buffers []MMAL_Buffer
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&buffers)))
	sliceHeader.Cap = int(pool.headers_num)
	sliceHeader.Len = int(pool.headers_num)
	sliceHeader.Data = uintptr(unsafe.Pointer(pool.header))
	return buffers
}

//export mmal_pool_callback
func mmal_pool_callback(pool *C.MMAL_POOL_T, buffer *C.MMAL_BUFFER_HEADER_T, userdata unsafe.Pointer) C.MMAL_BOOL_T {
	if cb, exists := pool_callback[pool]; exists {
		return mmal_to_bool(cb(pool, buffer, uintptr(userdata)))
	} else {
		// Empty buffer is available - queue it
		MMALQueuePut(pool.queue, buffer)
		return MMAL_BOOL_FALSE
	}
}
