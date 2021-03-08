// +build dispmanx

package surface

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Surface struct {
	sync.RWMutex

	x, y, w, h uint32
	layer      uint16
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSurface(x, y, w, h uint32) (*Surface, error) {
	this := new(Surface)

	// Check parameters
	if w == 0 || h == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewSurface")
	}

	// Return success
	return this, nil
}

func (this *Surface) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// TODO
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Surface) Origin() gopi.Point {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return gopi.Point{float32(this.x), float32(this.y)}

}

func (this *Surface) Size() gopi.Size {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return gopi.Size{float32(this.w), float32(this.h)}
}

func (this *Surface) Layer() uint16 {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.layer
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Surface) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<surface"
	str += fmt.Sprintf(" origin={%d,%d} size={%d,%d}", this.x, this.y, this.w, this.h)
	str += fmt.Sprint(" layer=", this.layer)

	return str + ">"
}
