// +build drm

package drm

import (
	"fmt"
	"os"
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
	Device     C.drmDevice
	DeviceNode uint
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MAX_DRM_DEVICES = 10
)

const (
	DRM_NODE_PRIMARY DeviceNode = C.DRM_NODE_PRIMARY
	DRM_NODE_CONTROL DeviceNode = C.DRM_NODE_CONTROL
	DRM_NODE_RENDER  DeviceNode = C.DRM_NODE_RENDER
	DRM_NODE_MAX                = C.DRM_NODE_MAX
)

////////////////////////////////////////////////////////////////////////////////
// DEVICE

func Devices() []*Device {
	devices := make([]C.drmDevicePtr, MAX_DRM_DEVICES)
	count := int(C.drmGetDevices2(0, &devices[0], C.int(MAX_DRM_DEVICES)))
	if count < 0 {
		return nil
	}
	result := make([]*Device, count)
	for i := 0; i < count; i++ {
		result[i] = (*Device)(unsafe.Pointer(devices[i]))
	}
	C.drmFreeDevices(&devices[0], C.int(count))
	return result
}

func (this *Device) Nodes() []string {
	var nodes []*C.char
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&nodes)))
	sliceHeader.Cap = int(DRM_NODE_MAX)
	sliceHeader.Len = int(DRM_NODE_MAX)
	sliceHeader.Data = uintptr(unsafe.Pointer(this.nodes))

	result := make([]string, len(nodes))
	for i, cstr := range nodes {
		if this.available_nodes&(1<<i) == 0 {
			continue
		} else if cstr != nil {
			result[i] = C.GoString(cstr)
		}
	}
	return result
}

func (this *Device) AvailableNode(node DeviceNode) string {
	node_ := int(node)
	if node_ >= DRM_NODE_MAX {
		return ""
	} else if this.available_nodes&(1<<node_) == 0 {
		return ""
	} else {
		return this.Nodes()[node_]
	}
}

func OpenDevice(device *Device, node DeviceNode) (*os.File, error) {
	if device == nil || node >= DeviceNode(DRM_NODE_MAX) {
		return nil, gopi.ErrBadParameter.WithPrefix("OpenDevice")
	} else if path := device.AvailableNode(node); path == "" {
		return nil, gopi.ErrBadParameter.WithPrefix("OpenDevice")
	} else {
		return OpenPath(path)
	}
}

func OpenPath(path string) (*os.File, error) {
	if fh, err := os.OpenFile(path, os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	} else {
		return fh, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Device) String() string {
	str := "<drm.device"
	nodes := this.Nodes()
	for n := DeviceNode(0); n < DRM_NODE_MAX; n++ {
		if name := nodes[n]; name != "" {
			str += fmt.Sprintf(" %v=%q", n, name)
		}
	}
	return str + ">"
}

func (n DeviceNode) String() string {
	switch n {
	case DRM_NODE_PRIMARY:
		return "DRM_NODE_PRIMARY"
	case DRM_NODE_CONTROL:
		return "DRM_NODE_CONTROL"
	case DRM_NODE_RENDER:
		return "DRM_NODE_RENDER"
	default:
		return "[?? Invalid DeviceNode value]"
	}
}
