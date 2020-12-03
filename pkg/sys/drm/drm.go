// +build drm

package drm

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/djthorpe/gopi/v2"
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
	ModeResources C.drmModeRes
	ModeConnector C.
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DRM_DEV = "/dev/dri/card"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - OPEN

func Device(bus uint) string {
	return fmt.Sprintf("%v%v", DRM_DEV, bus)
}

func Devices() []uint {
	matches, err := filepath.Glob(DRM_DEV + "*")
	if err != nil {
		return nil
	}
	var result []uint
	for _, match := range matches {
		if bus, err := strconv.ParseUint(strings.TrimPrefix(match, DRM_DEV), 10, 32); err == nil {
			result = append(result, uint(bus))
		}
	}
	return result
}

func OpenDevice(bus uint) (*os.File, error) {
	if file, err := os.OpenFile(Device(bus), os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET RESOURCES

func ModeGetResources(fd uintptr) (*ModeResources, error) {
	if res := C.drmModeGetResources(C.int(fd)); res == nil {
		return nil, gopi.ErrInternalAppError
	} else {
		return (*ModeResources)(res), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// ModeResources

func (this *ModeResources) FrameBuffers() []uint32 {
	
}

func (this *ModeResources) CRTCs() []uint32 {
	
}

func (this *ModeResources) Connectors() []uint32 {

}

func (this *ModeResources) Encoders() []uint32 {

}

func (this *ModeResources) Width() (uint32,uint32) {
	ctx := (*C.drmModeRes)(this)
	return uint32(ctx.min_width),uint32(ctx.max_width)
}

func (this *ModeResources) Height() (uint32,uint32) {
	ctx := (*C.drmModeRes)(this)
	return uint32(ctx.min_height),uint32(ctx.max_height)
}

func (this *ModeResources) String() string {
	str := "<drm.moderesources"
	if min,max := this.Width(); min > 0 && max > 0 {
		str += fmt.Sprintf(" width{min,max}={%v,%v}",min,max)
	}
	return str + ">"
}