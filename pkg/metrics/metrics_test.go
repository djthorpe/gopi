package metrics_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/metrics"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

type App struct {
	gopi.Unit
	gopi.Metrics
}

func Test_Metrics_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Metrics == nil {
			t.Error("nil metrics unit")
		} else if m, err := app.Metrics.NewMeasurement("name", "m1,m2,m3 float64"); err != nil {
			t.Error(err)
		} else if m.Name() != "name" {
			t.Error("Unexpected name", m.Name())
		} else if len(m.Metrics()) != 3 {
			t.Error("Unexpected metrics", m.Metrics())
		} else {
			t.Log("measurement=", m)
		}
	})
}

func Test_Metrics_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		fields := app.Metrics.NewFields("a", "b", "c")
		if len(fields) != 3 {
			t.Error("Unexpected number of fields")
		} else {
			t.Log(fields)
		}
	})
}
