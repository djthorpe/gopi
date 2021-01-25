// +build dvb

package dvb_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/media/dvb"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type ManagerApp struct {
	gopi.Unit
	gopi.DVBManager
}

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
		tuners := app.DVBManager.Tuners()
		for _, tuner := range tuners {
			t.Log(tuner)
		}
	})
}
