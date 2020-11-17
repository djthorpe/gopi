package metrics

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
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

func Test_Metrics_003(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		fields := app.Metrics.NewFields("a=true")
		if len(fields) != 1 {
			t.Error("Unexpected number of fields")
		} else if hasNilElement(fields) {
			t.Error("Field is nil")
		} else {
			t.Log(fields)
		}
	})
}

func Test_Metrics_004(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		fields := app.Metrics.NewFields("b=100", "c=200i")
		if len(fields) != 2 {
			t.Error("Unexpected number of fields")
		} else if hasNilElement(fields) {
			t.Error("Field is nil")
		} else if fields[0].Kind() != "uint64" {
			t.Error(fields[0], "Unexpected kind")
		} else if fields[1].Kind() != "int64" {
			t.Error(fields[1], "Unexpected kind")
		} else {
			t.Log(fields)
		}
	})
}

func Test_Metrics_005(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		fields := app.Metrics.NewFields("d=\"test\"")
		if len(fields) != 1 {
			t.Error("Unexpected number of fields")
		} else if hasNilElement(fields) {
			t.Error("Field is nil")
		} else if fields[0].Kind() != "string" {
			t.Error(fields[0], "Unexpected kind")
		} else {
			t.Log(fields)
		}
	})
}

func Test_Metrics_006(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if _, err := app.NewMeasurement("test", ""); err == nil {
			t.Error("Expected error for zero metrics")
		} else {
			t.Log("Expected error:", err)
		}
		if _, err := app.NewMeasurement("test", "a,a float64"); err == nil {
			t.Error("Expected error for duplicate metrics")
		} else {
			t.Log("Expected error:", err)
		}
		if _, err := app.NewMeasurement("test", "a,b float64", nil); err == nil {
			t.Error("Expected error for nil tag")
		} else {
			t.Log("Expected error:", err)
		}
		if _, err := app.NewMeasurement("test", "a,b float64", NewField("tag"), NewField("tag")); err == nil {
			t.Error("Expected error for duplicate tags")
		} else {
			t.Log("Expected error:", err)
		}
	})
}

func Test_Metrics_007(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if _, err := app.NewMeasurement("test", "a float64"); err != nil {
			t.Error(err)
		}
		if _, err := app.NewMeasurement("test", "a float64"); err == nil {
			t.Error("Expected error on duplicate measurement")
		} else {
			t.Log("Expected error:", err)
		}
	})
}

func Test_Metrics_008(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if len(app.Measurements()) != 0 {
			t.Error("Expected zero measurements")
		}
		app.NewMeasurement("test", "a float64")
		if len(app.Measurements()) != 1 {
			t.Error("Expected one measurement")
		}
		app.NewMeasurement("test2", "a float64")
		if len(app.Measurements()) != 2 {
			t.Error("Expected one measurement")
		}
	})
}

func Test_Metrics_009(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		app.NewMeasurement("test", "a bool")
		if err := app.Emit("test", true); err != nil {
			t.Error(err)
		}
		if err := app.Emit("test", false); err != nil {
			t.Error(err)
		}
		if err := app.Emit("test", 100); err == nil {
			t.Error("Expected error")
		} else {
			t.Log("Expected error", err)
		}
	})
}
