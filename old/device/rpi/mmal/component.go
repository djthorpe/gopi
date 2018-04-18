/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package mmal /* import "github.com/djthorpe/gopi/device/rpi/mmal" */

import (
	"fmt"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/mmal
    #cgo LDFLAGS:  -L/opt/vc/lib -lmmal -lmmal_components -lmmal_core
	#include <mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MMAL struct {
	// name of the component to create
	Name string
}

type Component struct {
	log     *util.LoggerDevice  // logger
	name    string              // name of the component
	handle  *C.MMAL_COMPONENT_T // the component handle
	control *Port
	input   []*Port
	output  []*Port
	clock   []*Port
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_COMPONENT_DEFAULT_VIDEO_DECODER   = "vc.ril.video_decode"
	MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER   = "vc.ril.video_encode"
	MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER  = "vc.ril.video_render"
	MMAL_COMPONENT_DEFAULT_IMAGE_DECODER   = "vc.ril.image_decode"
	MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER   = "vc.ril.image_encode"
	MMAL_COMPONENT_DEFAULT_CAMERA          = "vc.ril.camera"
	MMAL_COMPONENT_DEFAULT_VIDEO_CONVERTER = "vc.video_convert"
	MMAL_COMPONENT_DEFAULT_SPLITTER        = "vc.splitter"
	MMAL_COMPONENT_DEFAULT_SCHEDULER       = "vc.scheduler"
	MMAL_COMPONENT_DEFAULT_VIDEO_INJECTER  = "vc.video_inject"
	MMAL_COMPONENT_DEFAULT_VIDEO_SPLITTER  = "vc.ril.video_splitter"
	MMAL_COMPONENT_DEFAULT_AUDIO_DECODER   = "none"
	MMAL_COMPONENT_DEFAULT_AUDIO_RENDERER  = "vc.ril.audio_render"
	MMAL_COMPONENT_DEFAULT_MIRACAST        = "vc.miracast"
	MMAL_COMPONENT_DEFAULT_CLOCK           = "vc.clock"
	MMAL_COMPONENT_DEFAULT_CAMERA_INFO     = "vc.camera_info"
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new MMAL object, returns error if not possible
func (config MMAL) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<mmal.Component>Open name=%v", config.Name)

	// create new MMAL driver
	this := new(Component)
	this.log = log
	this.name = config.Name

	ret := status(C.mmal_component_create(C.CString(this.name), unsafe.Pointer(&this.handle)))
	if ret != MMAL_SUCCESS {
		return nil, ret.Error()
	}

	// Create ports
	this.control = &Port{handle: this.handle.control}
	this.input = mmalMakePorts(uint32(this.handle.input_num), this.handle.input)
	this.output = mmalMakePorts(uint32(this.handle.output_num), this.handle.output)
	this.clock = mmalMakePorts(uint32(this.handle.clock_num), this.handle.clock)

	return this, nil
}

// Private method to make a slice of ports
func mmalMakePorts(num uint32, ptr **C.MMAL_PORT_T) []*Port {
	if num == 0 {
		return nil
	}
	ports := make([]*Port, num)
	for i := uint32(0); i < num; i++ {
		ports[i] = &Port{handle: *ptr}
		ptr = (**C.MMAL_PORT_T)(unsafe.Pointer((uintptr(unsafe.Pointer(ptr)) + uintptr((unsafe.Sizeof(ptr))))))
	}
	return ports
}

// Close MMAL connection
func (this *Component) Close() error {
	this.log.Debug("<mmal.Component>Close name=%v", this.name)

	ret := status(C.mmal_component_destroy(unsafe.Pointer(this.handle)))
	if ret != MMAL_SUCCESS {
		return ret.Error()
	}

	return nil
}

func (this *Component) Name() string {
	return C.GoString(this.handle.name)
}

func (this *Component) Id() uint32 {
	return uint32(this.handle.id)
}

func (this *Component) String() string {
	return fmt.Sprintf("<mmal.Component>{ name=%v id=0x%08X enabled=%v control=%v input=%v output=%v clock=%v }", this.Name(), this.Id(), this.IsEnabled(), this.control, this.input, this.output, this.clock)
}

////////////////////////////////////////////////////////////////////////////////
// ENABLE AND DISABLE

func (this *Component) Enable() error {
	ret := status(C.mmal_component_enable(unsafe.Pointer(this.handle)))
	if ret != MMAL_SUCCESS {
		return ret.Error()
	}
	return nil
}

func (this *Component) Disable() error {
	ret := status(C.mmal_component_disable(unsafe.Pointer(this.handle)))
	if ret != MMAL_SUCCESS {
		return ret.Error()
	}
	return nil
}

func (this *Component) IsEnabled() bool {
	if uint32(this.handle.is_enabled) == 0 {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// PORTS

// Return an input port
func (this *Component) GetInput(i int) (*Port, error) {
	if i < 0 || i >= int(this.handle.input_num) {
		return nil, MMAL_EINVAL.Error()
	}
	return this.input[i], nil
}

// Return an output port
func (this *Component) GetOutput(i int) (*Port, error) {
	if i < 0 || i >= int(this.handle.output_num) {
		return nil, MMAL_EINVAL.Error()
	}
	return this.output[i], nil
}

// Return a clock port
func (this *Component) GetClock(i int) (*Port, error) {
	if i < 0 || i >= int(this.handle.clock_num) {
		return nil, MMAL_EINVAL.Error()
	}
	return this.clock[i], nil
}

// Return control port
func (this *Component) GetControl() *Port {
	return this.control
}

// Return the number of input ports
func (this *Component) NumberOfInputs() uint32 {
	return uint32(this.handle.input_num)
}

// Return the number of output ports
func (this *Component) NumberOfOutputs() uint32 {
	return uint32(this.handle.output_num)
}

// Return the number of clock ports
func (this *Component) NumberOfClocks() uint32 {
	return uint32(this.handle.clock_num)
}
