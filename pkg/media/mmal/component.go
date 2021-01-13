// +build mmal

package mmal

import (
	mmal "github.com/djthorpe/gopi/v3/pkg/sys/mmal"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type Component interface {
	Enable() error
	Disable() error
}

type ImageComponent interface {
	Component

	SetInputFormatJPEG() error
}

type VideoComponent interface {
	Component
}

type AudioComponent interface {
	Component
}

type CameraComponent interface {
	Component
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

type component struct {
	ctx *mmal.MMALComponent
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewComponent(ctx *mmal.MMALComponent) *component {
	return &component{ctx}
}

func (this *component) Dispose() error {
	var result error

	// If component enabled, then disable it
	if this.ctx.Enabled() {
		if err := this.ctx.Disable(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Free resources
	if err := this.ctx.Free(); err != nil {
		result = multierror.Append(result, err)
	}

	// Release resources
	this.ctx = nil

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *component) String() string {
	return this.ctx.String()
}

////////////////////////////////////////////////////////////////////////////////
// COMPONENT METHODS

func (this *component) Enable() error {
	return this.ctx.Enable()
}

func (this *component) Disable() error {
	return this.ctx.Disable()
}

////////////////////////////////////////////////////////////////////////////////
// IMAGE METHODS

func (this *component) SetInputFormatJPEG() error {
	port := this.ctx.InputPorts()[0]
	port.Format().SetEncoding(mmal.MMAL_ENCODING_JPEG)
	return port.FormatCommit()
}
