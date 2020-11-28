package influxdb_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	influxdb "github.com/djthorpe/gopi/v3/pkg/db/influxdb"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type WriterApp struct {
	gopi.Unit
	*influxdb.Writer
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

func Test_Writer_002(t *testing.T) {
	tool.Test(t, []string{"-influxdb.url=rpi4b"}, new(WriterApp), func(app *WriterApp) {
		if delta, err := app.Writer.Ping(); err != nil {
			t.Error(err)
		} else {
			t.Log("Ping delta=", delta)
		}
	})
}
