// +build dvb

package dvb_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/media/dvb"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type ManagerApp struct {
	gopi.Unit
	gopi.DVBManager
}

var (
	PATH_TUNE_PARAMS = "/usr/share/dvb/dvb-c/de-Kabel_Deutschland-Hannover"
)

func Test_Manager_001(t *testing.T) {
	tool.Test(t, nil, new(ManagerApp), func(app *ManagerApp) {
		if app.DVBManager == nil {
			t.Error("manager is nil")
		} else {
			t.Log(app.DVBManager)
		}
	})
}

func Test_Manager_002(t *testing.T) {
	tool.Test(t, nil, new(ManagerApp), func(app *ManagerApp) {
		if _, err := os.Stat(PATH_TUNE_PARAMS); os.IsNotExist(err) {
			t.Skip("Skipping test, no file")
		}
		fh, err := os.Open(PATH_TUNE_PARAMS)
		if err != nil {
			t.Fatal(err)
		}
		defer fh.Close()
		if params, err := app.ParseTunerParams(fh); err != nil {
			t.Error(err)
		} else {
			t.Log(params)
		}
	})
}

func Test_Manager_003(t *testing.T) {
	tool.Test(t, nil, new(ManagerApp), func(app *ManagerApp) {
		devices := app.DVBManager.Tuners()
		if len(devices) == 0 {
			t.Skip("Skipping test, no device")
		}
		if _, err := os.Stat(PATH_TUNE_PARAMS); os.IsNotExist(err) {
			t.Skip("Skipping test, no file")
		}
		fh, err := os.Open(PATH_TUNE_PARAMS)
		if err != nil {
			t.Fatal(err)
		}
		defer fh.Close()
		params, err := app.ParseTunerParams(fh)
		if err != nil {
			t.Error(err)
		}
		for _, param := range params {
			t.Log("Tuning", param.Name())
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()
			if err := app.DVBManager.Tune(ctx, devices[0], param); errors.Is(err, context.DeadlineExceeded) {
				t.Log("  Tune Timeout")
			} else if err != nil {
				t.Error(err)
			} else {
				t.Log("  Tune OK")
			}
		}
	})
}
