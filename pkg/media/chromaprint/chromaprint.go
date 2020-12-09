// +build chromaprint

package chromaprint

import (
	"strconv"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	chromaprint "github.com/djthorpe/gopi/v3/pkg/sys/chromaprint"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Logger
	sync.Mutex

	ctx []*chromaprint.Context
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func (this *Manager) New(gopi.Config) error {

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	for _, ctx := range this.ctx {
		if ctx != nil {
			ctx.Free()
		}
	}

	// Release resources
	this.ctx = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<manager.chromaprint"
	if v := chromaprint.Version(); v != "" {
		str += " version=" + strconv.Quote(v)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) NewStream(rate, channels int) (*chromaprint.Context, error) {
	if ctx := chromaprint.NewChromaprint(chromaprint.ALGORITHM_DEFAULT); ctx == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("NewStream")
	} else if err := ctx.Start(rate, channels); err != nil {
		ctx.Free()
		return nil, err
	} else {
		this.ctx = append(this.ctx, ctx)
		return ctx, nil
	}
}

func (this *Manager) Close(ctx *chromaprint.Context) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	for i := range this.ctx {
		if ctx == this.ctx[i] {
			ctx.Free()
			this.ctx[i] = nil
			return nil
		}
	}

	return gopi.ErrNotFound.WithPrefix("Close")
}
