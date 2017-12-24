/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host
    #cgo LDFLAGS:  -L/opt/vc/lib -lbcm_host
	#include "vc_dispmanx.h"
*/
import "C"

////////////////////////////////////////////////////////////////////////////////

type (
	dxUpdateHandle   uint32
	dxUpdatePriority int32
)

const (
	DX_UPDATE_NONE dxUpdateHandle = 0
)

////////////////////////////////////////////////////////////////////////////////

func (this *DXDisplay) UpdateBegin() (dxUpdateHandle, error) {
	handle := updateStart(dxUpdatePriority(0))
	if handle == DX_UPDATE_NONE {
		return handle, this.log.Error("dxUpdateStart failed")
	}
	return handle, nil
}

func (this *DXDisplay) UpdateSubmit(handle dxUpdateHandle) error {
	if updateSubmitSync(handle) != true {
		return this.log.Error("dxUpdateSubmitSync failed")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods - updates

func updateStart(priority dxUpdatePriority) dxUpdateHandle {
	return dxUpdateHandle(C.vc_dispmanx_update_start(C.int32_t(priority)))
}

func updateSubmitSync(handle dxUpdateHandle) bool {
	return C.vc_dispmanx_update_submit_sync(C.DISPMANX_UPDATE_HANDLE_T(handle)) == DX_SUCCESS
}
