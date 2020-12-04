package csv_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
	_ "github.com/djthorpe/gopi/v3/pkg/metrics"
)

type WriterApp struct {
	gopi.Unit
	gopi.Metrics
	gopi.MetricWriter
	gopi.Platform
}

func (this *WriterApp) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
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
	tempdir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)
	tool.Test(t, []string{"-csv.path", tempdir}, new(WriterApp), func(app *WriterApp) {
		// Make new measurement
		if _, err := app.Metrics.NewMeasurement("test", "m1,m5,m15 float64", app.Metrics.HostTag()); err != nil {
			t.Error(err)
		}

		// Write metrics
		for i := 0; i < 10; i++ {
			t.Log("Writing metric", i)
			l1, l5, l15 := app.Platform.LoadAverages()
			if err := app.Metrics.Emit("test", float64(l1), float64(l5), float64(l15)); err != nil {
				t.Error(err)
			}
		}

		// TODO: Flush (wait for empty publisher)
		time.Sleep(time.Second)

		// Read and parse CSV file
		if data, err := ioutil.ReadFile(filepath.Join(tempdir, "test.csv")); err != nil {
			t.Error(err)
		} else {
			t.Log("data=", string(data))
		}
	})
}

func Test_Writer_003(t *testing.T) {
	tempdir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	tool.Test(t, []string{"-csv.path", tempdir}, new(WriterApp), func(app *WriterApp) {
		// Make new measurement
		if _, err := app.Metrics.NewMeasurement("test", "m1,m5,m15 float64", app.Metrics.HostTag()); err != nil {
			t.Error(err)
		}

		// Write metrics
		for i := 0; i < 10; i++ {
			t.Log("Writing metric", i)
			l1, l5, l15 := app.Platform.LoadAverages()
			if err := app.Metrics.EmitTS("test", time.Now(), float64(l1), float64(l5), float64(l15)); err != nil {
				t.Error(err)
			}
		}

		// TODO: Flush (wait for empty publisher)
		time.Sleep(time.Second)

		// Read and parse CSV file
		if data, err := ioutil.ReadFile(filepath.Join(tempdir, "test.csv")); err != nil {
			t.Error(err)
		} else {
			t.Log("data=", string(data))
		}
	})

	tool.Test(t, []string{"-csv.path", tempdir}, new(WriterApp), func(app *WriterApp) {
		// Make new measurement
		if _, err := app.Metrics.NewMeasurement("test", "m1,m5,m15 float64", app.Metrics.HostTag()); err != nil {
			t.Error(err)
		}

		// Write metrics
		for i := 0; i < 10; i++ {
			t.Log("Writing metric", i)
			l1, l5, l15 := app.Platform.LoadAverages()
			if err := app.Metrics.EmitTS("test", time.Now(), float64(l1), float64(l5), float64(l15)); err != nil {
				t.Error(err)
			}
		}

		// TODO: Flush (wait for empty publisher)
		time.Sleep(time.Second)

		// Read and parse CSV file
		if data, err := ioutil.ReadFile(filepath.Join(tempdir, "test.csv")); err != nil {
			t.Error(err)
		} else {
			t.Log("data=", string(data))
		}
	})
}
