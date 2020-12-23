// +build drm

package drm

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libdrm
#include <xf86drm.h>
#include <xf86drmMode.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	ModeConnector C.drmModeConnector
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func GetConnector(fd uintptr, id uint32) (*ModeConnector, error) {
	if conn := C.drmModeGetConnector(C.int(fd), C.uint32_t(id)); conn == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("GetConnector")
	} else {
		return (*ModeConnector)(unsafe.Pointer(conn)), nil
	}
}

func (this *ModeConnector) Free() {
	ctx := (*C.drmModeConnector)(this)
	C.drmModeFreeConnector(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *ModeConnector) Id() uint32 {
	ctx := (*C.drmModeConnector)(this)
	return uint32(ctx.connector_id)
}

func (this *ModeConnector) Status() ModeConnection {
	ctx := (*C.drmModeConnector)(this)
	return ModeConnection(ctx.connection)
}

func (this *ModeConnector) Type() ConnectorType {
	ctx := (*C.drmModeConnector)(this)
	return ConnectorType(ctx.connector_type)
}

func (this *ModeConnector) Dimensions() (uint32, uint32) {
	ctx := (*C.drmModeConnector)(this)
	return uint32(ctx.mmWidth), uint32(ctx.mmHeight)
}

func (this *ModeConnector) Modes() []ModeInfo {
	var result []ModeInfo

	// Make fake slice
	ctx := (*C.drmModeConnector)(this)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.count_modes)
	sliceHeader.Len = int(ctx.count_modes)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.modes))

	return result
}

func (this *ModeConnector) Encoder() uint32 {
	ctx := (*C.drmModeConnector)(this)
	return uint32(ctx.encoder_id)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *ModeConnector) String() string {
	str := "<drm.connector"
	str += " id=" + fmt.Sprint(this.Id())
	if c := this.Status(); c != ModeConnectionNone {
		str += " status=" + fmt.Sprint(c)
	}
	if c := this.Type(); c != DRM_MODE_CONNECTOR_Unknown {
		str += " type=" + fmt.Sprint(c)
	}
	if enc := this.Encoder(); enc != 0 {
		str += " encoder=" + fmt.Sprint(enc)
	}
	if w, h := this.Dimensions(); w > 0 && h > 0 {
		str += fmt.Sprintf(" dimensions={%vmm,%vmm}", w, h)
	}
	str += fmt.Sprintf(" modes=%v", this.Modes())
	return str + ">"
}
