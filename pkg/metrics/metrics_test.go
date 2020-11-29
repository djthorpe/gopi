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
			t.Error("NewMeasurement", err)
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
		a := app.Metrics.Field("a")
		if a == nil {
			t.Error("Unexpected nil returned")
		} else {
			t.Log(a)
		}
	})
}

func Test_Metrics_003(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if field := app.Metrics.Field("a", true); field == nil {
			t.Error("Field is nil")
		} else if k := field.Name(); k != "a" {
			t.Error("Unexpected field name")
		} else if k := field.Kind(); k != "bool" {
			t.Error("Unexpected field kind")
		} else if v := field.Value().(bool); v != true {
			t.Error("Field is not true")
		} else {
			t.Log(field)
		}
	})
}

func Test_Metrics_004(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		field := app.Metrics.Field("b", int64(100))
		if field == nil {
			t.Error("Field is nil")
		} else if field.Kind() != "int64" {
			t.Error(field, "Unexpected kind")
		} else if field.Value().(int64) != int64(100) {
			t.Error(field, "Unexpected value")
		} else {
			t.Log(field)
		}
	})
}

func Test_Metrics_005(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		field := app.Metrics.Field("d")
		if field == nil {
			t.Error("Field is nil")
		} else if field.Kind() != "nil" {
			t.Error(field, "Unexpected kind")
		}
		if err := field.SetValue("test"); err != nil {
			t.Error(err)
		} else if field.Kind() != "string" {
			t.Error(field, "Unexpected kind")
		} else if field.Value().(string) != "test" {
			t.Error(field, "Unexpected value")
		}

		if err := field.SetValue(nil); err != nil {
			t.Error(err)
		} else if field.Kind() != "string" {
			t.Error(field, "Unexpected kind")
		} else if field.IsNil() == false {
			t.Error(field, "Unexpected IsNil value")
		} else if field.Value().(string) != "" {
			t.Error(field, "Unexpected value")
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
