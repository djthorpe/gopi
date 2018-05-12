/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi_test

import (
	"testing"
	"time"

	// Import frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/hw/mock"
	_ "github.com/djthorpe/gopi/sys/logger"
)

// Create an app with metrics module
func TestMetrics_000(t *testing.T) {
	config := gopi.NewAppConfig("metrics")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else {
		t.Log("app", app)
	}
}

// Cast to gopi.Metrics
func TestMetrics_001(t *testing.T) {
	config := gopi.NewAppConfig("metrics")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if metrics, ok := app.ModuleInstance("metrics").(gopi.Metrics); ok == false {
		_ = app.ModuleInstance("metrics").(gopi.Metrics)
		t.Fatal("Unable to cast to metrics interface")
	} else if metrics == nil {
		t.Fatal("metrics == nil")
	}
}

// Retrieve information about load average, etc
func TestMetrics_002(t *testing.T) {
	config := gopi.NewAppConfig("metrics")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if metrics, ok := app.ModuleInstance("metrics").(gopi.Metrics); ok == false {
		_ = app.ModuleInstance("metrics").(gopi.Metrics)
		t.Fatal("Unable to cast to metrics interface")
	} else {
		uptime_host := metrics.UptimeHost()
		if uptime_host == 0 {
			t.Error("uptime_host == 0")
		}
		time.Sleep(time.Second * 1)
		if uptime_host >= metrics.UptimeHost() {
			t.Error("uptime_host failure")
		}
	}
}

// Retrieve information about load average, etc
func TestMetrics_003(t *testing.T) {
	config := gopi.NewAppConfig("metrics")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if metrics, ok := app.ModuleInstance("metrics").(gopi.Metrics); ok == false {
		_ = app.ModuleInstance("metrics").(gopi.Metrics)
		t.Fatal("Unable to cast to metrics interface")
	} else {
		uptime_app := metrics.UptimeApp()
		if uptime_app == 0 {
			t.Error("uptime_app == 0")
		}
		time.Sleep(time.Second * 1)
		if uptime_app >= metrics.UptimeApp() {
			t.Error("uptime_app failure")
		}
	}
}

func TestMetrics_004(t *testing.T) {
	config := gopi.NewAppConfig("metrics")
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if metrics, ok := app.ModuleInstance("metrics").(gopi.Metrics); ok == false {
		_ = app.ModuleInstance("metrics").(gopi.Metrics)
		t.Fatal("Unable to cast to metrics interface")
	} else {
		uptime_app := metrics.UptimeApp()
		uptime_host := metrics.UptimeHost()
		if uptime_app <= uptime_host {
			t.Error("uptime_app > uptime_host")
		}
	}
}

func TestMetrics_005(t *testing.T) {
	config := gopi.NewAppConfig("metrics")
	config.Debug = true
	config.Verbose = true
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Fatal(err)
	} else if metrics, ok := app.ModuleInstance("metrics").(gopi.Metrics); ok == false {
		_ = app.ModuleInstance("metrics").(gopi.Metrics)
		t.Fatal("Unable to cast to metrics interface")
	} else if counter, err := metrics.NewCounter(gopi.METRIC_TYPE_NONE, gopi.METRIC_RATE_SECOND, "Counter"); err != nil {
		t.Error(err)
	} else {
		for i := 0; i < 100; i++ {
			// Increment the counter by one
			counter <- 1
			// wait for a second
			time.Sleep(time.Second)
			// Get metric
			if metric := metrics.Metric(counter); metric == nil {
				t.Error("metric==nil")
			} else if metric.Mean != 1 {
				t.Error("mean != 1 per second")
			} else if metric.Total != 60 {
				t.Error("total != 60 per hour")
			} else {
				t.Logf("Interation %v mean=%v total=%v", i, metric.Mean, metric.Total)
			}
		}

		// Close the application
		app.Close()
	}
}
