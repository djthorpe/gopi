// +build linux

package input_test

import (
	"context"
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/file"
)

type InputApp struct {
	gopi.Unit
	gopi.Logger
	gopi.InputManager
}

func (this *InputApp) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func Test_InputManager_001(t *testing.T) {
	tool.Test(t, nil, new(InputApp), func(app *InputApp) {
		app.Require(app.Logger, app.InputManager)
	})
}
