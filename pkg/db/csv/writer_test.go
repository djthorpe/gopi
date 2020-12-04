package csv_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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
	tempdir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempdir)
	tool.Test(t, []string{"-csv.path", tempdir}, new(WriterApp), func(app *WriterApp) {
		_, err := app.Metrics.NewMeasurement("test", "m1,m5,m15 float64", app.Metrics.HostTag())
		if err != nil {
			t.Error(err)
		}

		// Write metrics
		for i := 0; i < 10; i++ {
			if err := app.Metrics.Emit("test", float64(1), float64(1), float64(1)); err != nil {
				t.Error(err)
			}
		}

		// TODO: Flush (wait for empty publisher)
		time.Sleep(1 * time.Second)

		// Read and parse CSV file
		if data, err := ioutil.ReadFile(filepath.Join(tempdir, "test.csv")); err != nil {
			t.Error(err)
		} else {
			t.Log("data=", string(data))
		}
	})
}
