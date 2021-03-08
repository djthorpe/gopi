// +build dispmanx

package surface

import (
	"fmt"
	"sync"

	"github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Context struct {
	sync.RWMutex

	dx.Display
	dx.Update
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewContext(display dx.Display, priority int32) (*Context, error) {
	this := new(Context)
	if update, err := dx.UpdateStart(priority); err != nil {
		return nil, err
	} else {
		this.Display = display
		this.Update = update
	}

	// Return success
	return this, nil
}

func (this *Context) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Can't dispose twice
	if this.Update == 0 {
		return gopi.ErrOutOfOrder
	}

	// Submit sync
	var result error
	if err := dx.UpdateSubmitSync(this.Update); err != nil {
		result = multierror.Append(result, err)
	} else {
		this.Update = 0
	}

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Context) Valid() bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.Update != 0
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Context) String() string {
	str := "<context"
	if handle := this.Update; handle != 0 {
		str += fmt.Sprintf(" dispmanx=0x%08X", this.Update)
	}
	return str + ">"
}
