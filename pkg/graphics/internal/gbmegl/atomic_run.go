// +build egl

package gbmegl

import (
	"context"

	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

////////////////////////////////////////////////////////////////////////////////
// RUN LOOP

func (this *EGL) SwapBuffers() error {
	var result error
	for _, surface := range this.surfaces {
		if err := surface.MakeCurrent(this.display); err != nil {
			result = multierror.Append(err)
		} else if err := surface.Draw(); err != nil {
			result = multierror.Append(err)
		}
		if err := surface.SwapBuffers(this.display); err != nil {
			result = multierror.Append(err)
		}
	}
	return result
}

func (this *EGL) Run(ctx context.Context) error {
	for _ = range ctx.Done() {
		if err := this.SwapBuffers(); err != nil {
			return err
		}
		// TODO here
	}
	return ctx.Err()
}
