package gbm

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: gbm
#include <gbm.h>
#include <errno.h>

int _errno() { return errno; }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	GBMBufferFlags C.enum_gbm_bo_flags
	GBMBufferType  C.uint32_t
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Buffer is going to be presented to the screen using an API such as KMS
	GBM_BO_USE_SCANOUT GBMBufferFlags = C.GBM_BO_USE_SCANOUT

	// Buffer is going to be used as cursor
	GBM_BO_USE_CURSOR GBMBufferFlags = C.GBM_BO_USE_CURSOR

	// Buffer is to be used for rendering
	GBM_BO_USE_RENDERING GBMBufferFlags = C.GBM_BO_USE_RENDERING

	// Buffer can be used for gbm_bo_write
	GBM_BO_USE_WRITE GBMBufferFlags = C.GBM_BO_USE_WRITE

	// Buffer is linear, i.e. not tiled.
	GBM_BO_USE_LINEAR GBMBufferFlags = C.GBM_BO_USE_LINEAR
)

const (
	GBM_BO_IMPORT_WL_BUFFER   GBMBufferType = C.GBM_BO_IMPORT_WL_BUFFER
	GBM_BO_IMPORT_EGL_IMAGE   GBMBufferType = C.GBM_BO_IMPORT_WL_BUFFER
	GBM_BO_IMPORT_FD          GBMBufferType = C.GBM_BO_IMPORT_FD
	GBM_BO_IMPORT_FD_MODIFIER GBMBufferType = C.GBM_BO_IMPORT_FD_MODIFIER
)
