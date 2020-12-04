package csv_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
	_ "github.com/djthorpe/gopi/v3/pkg/metrics"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type WriterApp struct {
	gopi.Unit
	gopi.Metrics
	gopi.MetricWriter
}

func Test_Writer_001(t *testing.T) {
	tool.Test(t, nil, new(WriterApp), func(app *WriterApp) {
		if app.MetricWriter == nil {
			t.Error("MetricWriter is nil")
		} else {
			t.Log(app.MetricWriter)
		}
	})
}
func Test_Writer_002(t *testing.T) {
	tool.Test(t, nil, new(WriterApp), func(app *WriterApp) {
		m, err := app.Metrics.NewMeasurement("test", "test string")
		if err != nil {
			t.Error(err)
		}
		m.Set("test", "hello, world")
		if err := app.MetricWriter.Write(m); err != nil {
			t.Error(err)
		}
	})
}
