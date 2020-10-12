package display

import (
	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/hw/platform"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Displays struct {
	gopi.Unit
	*platform.Platform

	displays map[uint16]gopi.Display
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// New is called to initialize
func (this *Displays) New(gopi.Config) error {
	this.displays = make(map[uint16]gopi.Display)
	return nil
}

// Dispose is called to close
func (this *Displays) Dispose() error {
	var result error
	for k := range this.displays {
		if err := this.Close(k); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

// Open returns a gopi.Display object based on id
func (this *Displays) Open(id uint16) (gopi.Display, error) {
	display := new(display)
	if display_, exists := this.displays[id]; exists {
		return display_, nil
	} else if err := display.new(id); err != nil {
		return nil, err
	} else {
		this.displays[id] = display
	}

	// Return success
	return display, nil
}

// Close disposes of an open display
func (this *Displays) Close(id uint16) error {
	if display_, exists := this.displays[id]; exists == false {
		return gopi.ErrBadParameter
	} else {
		delete(this.displays, id)
		return display_.(*display).close()
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Displays) String() string {
	str := "<displays"
	return str + ">"
}
