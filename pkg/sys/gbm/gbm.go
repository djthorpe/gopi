// +build gbm

package gbm

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: gbm
#include <gbm.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	GBMDevice C.struct_gbm_device
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GBM_DEVICE_PATH = "/dev/dri"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func GBMDevicePath(node string) string {
	return filepath.Join(GBM_DEVICE_PATH, node)
}

func GBMDevices() []string {
	files, err := ioutil.ReadDir(GBM_DEVICE_PATH)
	if err != nil {
		return nil
	}
	result := []string{}
	for _, file := range files {
		if file.Mode()&os.ModeDevice == 0 {
			continue
		}
		result = append(result, file.Name())
	}
	return result
}

func OpenDevice(node string) (*os.File, error) {
	if _, err := os.Stat(GBMDevicePath(node)); os.IsNotExist(err) {
		return nil, gopi.ErrNotFound.WithPrefix(node)
	} else if fh, err := os.Open(GBMDevicePath(node)); err != nil {
		return nil, err
	} else {
		return fh, nil
	}
}

func GBMCreateDevice(fd uintptr) *GBMDevice {
	if dev := C.gbm_create_device(C.int(fd)); dev == nil {
		return nil
	} else {
		return (*GBMDevice)(dev)
	}
}

func (this *GBMDevice) Free() {
	ctx := (*C.struct_gbm_device)(this)
	C.gbm_device_destroy(ctx)
}

func (this *GBMDevice) Name() string {
	ctx := (*C.struct_gbm_device)(this)
	if name := C.gbm_device_get_backend_name(ctx); name == nil {
		return ""
	} else {
		return C.GoString(name)
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *GBMDevice) String() string {
	str := "<gbm.device"
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	return str + ">"
}
