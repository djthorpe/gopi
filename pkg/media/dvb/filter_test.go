// +build dvb

package dvb_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/file"
	dvb "github.com/djthorpe/gopi/v3/pkg/media/dvb"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type FilterApp struct {
	gopi.Unit
	gopi.DVBManager
}

func Test_Filter_001(t *testing.T) {
	tool.Test(t, nil, new(FilterApp), func(app *FilterApp) {
		tuners := app.DVBManager.Tuners()
		if len(tuners) == 0 {
			t.Skip("Skipping test, no device")
		}
		for _, tuner := range tuners {
			if f, err := dvb.NewSectionFilter(tuner.(*dvb.Tuner), 0xFFFF, 0xFF); err != nil {
				t.Error(err)
			} else {
				t.Log(f)
				f.Dispose()
			}
		}
	})
}
