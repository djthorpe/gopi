/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"errors"
	"sync"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/mmal
    #cgo LDFLAGS:  -L/opt/vc/lib -lbcm_host
	#include "mmal.h"
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type MMAL struct {
	Device gopi.HardwareDriver
	Master uint
}

type MMDriver struct {
	log     *util.LoggerDevice // logger
	memlock sync.Mutex
	master  uint
	mem8    []uint8  // access I2C registers as bytes
	mem32   []uint32 // access I2C registers as uint32
}

type MMComponent string

type mmComponent *C.MMAL_COMPONENT_T

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_COMPONENT_DEFAULT_VIDEO_DECODER   MMComponent = "vc.ril.video_decode"
	MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER   MMComponent = "vc.ril.video_encode"
	MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER  MMComponent = "vc.ril.video_render"
	MMAL_COMPONENT_DEFAULT_IMAGE_DECODER   MMComponent = "vc.ril.image_decode"
	MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER   MMComponent = "vc.ril.image_encode"
	MMAL_COMPONENT_DEFAULT_CAMERA          MMComponent = "vc.ril.camera"
	MMAL_COMPONENT_DEFAULT_VIDEO_CONVERTER MMComponent = "vc.video_convert"
	MMAL_COMPONENT_DEFAULT_SPLITTER        MMComponent = "vc.splitter"
	MMAL_COMPONENT_DEFAULT_SCHEDULER       MMComponent = "vc.scheduler"
	MMAL_COMPONENT_DEFAULT_VIDEO_INJECTER  MMComponent = "vc.video_inject"
	MMAL_COMPONENT_DEFAULT_VIDEO_SPLITTER  MMComponent = "vc.ril.video_splitter"
	MMAL_COMPONENT_DEFAULT_AUDIO_DECODER   MMComponent = "none"
	MMAL_COMPONENT_DEFAULT_AUDIO_RENDERER  MMComponent = "vc.ril.audio_render"
	MMAL_COMPONENT_DEFAULT_MIRACAST        MMComponent = "vc.miracast"
	MMAL_COMPONENT_DEFAULT_CLOCK           MMComponent = "vc.clock"
	MMAL_COMPONENT_DEFAULT_CAMERA_INFO     MMComponent = "vc.camera_info"
)


////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new MMAL object, returns error if not possible
func (config MMAL) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<rpi.MMAL>Open")

	// create new GPIO driver
	this := new(MMDriver)

	// Set logging & device
	this.log = log

	// success
	return this, nil
}

// Close MMAL connection
func (this *MMDriver) Close() error {
	this.log.Debug("<rpi.MMAL>Close")

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// CREATE AND DESTROY COMPONENT

func (this *MMDriver) CreateComponent(component MMComponent) (*mmComponent,error) {
	return nil,errors.New("NOT IMPLEMENTED")
}

func (this *MMDriver) DestroyComponent(component *mmComponent) error {
	return errors.New("NOT IMPLEMENTED")
}




