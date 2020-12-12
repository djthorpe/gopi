package influxdb_test

import (
	"fmt"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	influxdb "github.com/djthorpe/gopi/v3/pkg/db/influxdb"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
	_ "github.com/djthorpe/gopi/v3/pkg/metrics"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type WriterApp struct {
	gopi.Unit
	gopi.Metrics
	*influxdb.Writer
}

type WriterServerApp struct {
	gopi.Unit
	gopi.Metrics
	gopi.Platform
	*influxdb.Writer
	*MockWriter
}

func Test_Writer_001(t *testing.T) {
	tool.Test(t, nil, new(WriterApp), func(app *WriterApp) {
		if app.Writer == nil {
			t.Error("writer is nil")
		} else if ep := app.Writer.Endpoint(); ep == nil {
			t.Error("Unexpected nil")
		} else if ep.Scheme != influxdb.DefaultScheme && ep.Host != "localhost" {
			t.Error("Unexpected endpoint", ep)
		} else {
			t.Log("endpoint=", ep)
		}
	})
}
func Test_Writer_002(t *testing.T) {
	tool.Test(t, []string{"-influxdb.url=host/database"}, new(WriterApp), func(app *WriterApp) {
		if app.Writer == nil {
			t.Error("writer is nil")
		} else if ep := app.Writer.Endpoint(); ep == nil {
			t.Error("Unexpected nil")
		} else if ep.Scheme != influxdb.DefaultScheme && ep.Host != "host" {
			t.Error("Unexpected endpoint", ep)
		} else if db := app.Writer.Database(); db != "database" {
			t.Error("Unexpected database", db)
		} else {
			t.Log("endpoint=", ep, " database=", db)
		}
	})
}
func Test_Writer_003(t *testing.T) {

	tool.Test(t, []string{"-influxdb.url=rpi4b"}, new(WriterApp), func(app *WriterApp) {
		if ep := app.Writer.Endpoint(); ep == nil {
			t.Error("Unexpected nil")
		} else if ep.Scheme != influxdb.DefaultScheme && ep.Host != "rpi4b" {
			t.Error("Unexpected endpoint", ep)
		} else if db := app.Writer.Database(); db != "" {
			t.Error("Unexpected database", db)
		} else {
			t.Log("endpoint=", ep)
		}
	})
}
func Test_Writer_004(t *testing.T) {

	tool.Test(t, []string{"-influxdb.url=rpi4b:9999"}, new(WriterApp), func(app *WriterApp) {
		if ep := app.Writer.Endpoint(); ep == nil {
			t.Error("Unexpected nil")
		} else if ep.Scheme != influxdb.DefaultScheme && ep.Host != "rpi4b" {
			t.Error("Unexpected endpoint", ep)
		} else if ep.Port() != "9999" {
			t.Error("Unexpected port", ep.Port())
		} else {
			t.Log("endpoint=", ep)
		}
	})
}
func Test_Writer_005(t *testing.T) {
	tool.Test(t, []string{"-influxdb.url=rpi4b:9999/metrics"}, new(WriterApp), func(app *WriterApp) {
		if ep := app.Writer.Endpoint(); ep == nil {
			t.Error("Unexpected nil")
		} else if ep.Scheme != influxdb.DefaultScheme && ep.Host != "rpi4b" {
			t.Error("Unexpected endpoint", ep)
		} else if ep.Port() != "9999" {
			t.Error("Unexpected port", ep.Port())
		} else if db := app.Writer.Database(); db != "metrics" {
			t.Error("Unexpected database", db)
		} else {
			t.Log("endpoint=", ep, " database=", db)
		}
	})
}

func Test_Writer_006(t *testing.T) {
	tool.Test(t, nil, new(WriterApp), func(app *WriterApp) {
		listen := fmt.Sprint("localhost:", influxdb.DefaultPort)
		if server, err := NewMockServer(t, listen); err != nil {
			t.Error(err)
		} else if delta, err := app.Writer.Ping(); err != nil {
			t.Error(err)
		} else if err := server.Close(); err != nil {
			t.Error(err)
		} else {
			t.Log("Ping delta=", delta)
		}
	})
}

func Test_Writer_007(t *testing.T) {
	tool.Test(t, nil, new(WriterApp), func(app *WriterApp) {
		listen := fmt.Sprint("localhost:", influxdb.DefaultPort)
		server, err := NewMockServer(t, listen)
		if err != nil {
			t.Error(err)
		}
		m, err := app.Metrics.NewMeasurement("test", "test string")
		if err != nil {
			t.Error(err)
		}
		m.Set("test", "hello, world")
		if err := app.Writer.Write(m); err != nil {
			t.Error(err)
		}
		if err := server.Close(); err != nil {
			t.Error(err)
		}
	})
}

func Test_Writer_008(t *testing.T) {
	tool.Test(t, nil, new(WriterServerApp), func(app *WriterServerApp) {
		if _, err := app.Metrics.NewMeasurement("loadavg", "l1,l5,l15 float64"); err != nil {
			t.Error(err)
		}

		// Emit metrics and have the mockwriter write to the database
		for i := 0; i < 10; i++ {
			l1, l5, l15 := app.Platform.LoadAverages()
			if err := app.Metrics.Emit("loadavg", l1, l5, l15); err != nil {
				t.Error(err)
			}
			time.Sleep(100 * time.Millisecond)
		}

		// Emit metrics and have the mockwriter write to the database
		for i := 0; i < 10; i++ {
			l1, l5, l15 := app.Platform.LoadAverages()
			if err := app.Metrics.EmitTS("loadavg", time.Now(), l1, l5, l15); err != nil {
				t.Error(err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	})
}

func Test_Writer_009(t *testing.T) {
	tool.Test(t, []string{"-influxdb.url=https://rpi4b/metrics"}, new(WriterApp), func(app *WriterApp) {
		if ep := app.Writer.Endpoint(); ep == nil {
			t.Error("Unexpected nil")
		} else if ep.String() != "https://rpi4b:8086/" {
			t.Error("Unexpected endpoint", ep)
		} else {
			t.Log("endpoint=", ep)
		}
	})
}

func Test_Writer_010(t *testing.T) {
	tool.Test(t, []string{"-influxdb.url=http://user:pass@rpi4:99/metrics"}, new(WriterApp), func(app *WriterApp) {
		if ep := app.Writer.Endpoint(); ep == nil {
			t.Error("Unexpected nil")
		} else if ep.User.String() != "" {
			t.Error("Unexpected username/password", ep)
		} else {
			t.Log("endpoint=", ep)
		}
	})
}
