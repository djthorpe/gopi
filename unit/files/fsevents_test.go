// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
	_ "github.com/djthorpe/gopi/v2/unit/files"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_FSEvents_000(t *testing.T) {
	t.Log("Test_FSEvents_000")
}

func Test_FSEvents_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_FSEvents_001, []string{"-debug"}, "gopi/fsevents"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_FSEvents_001(app gopi.App, t *testing.T) {
	fsevents := app.UnitInstance("gopi/fsevents").(gopi.FSEvents)
	t.Log(fsevents)
}

func Test_FSEvents_002(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_FSEvents_002, []string{"-debug"}, "gopi/fsevents"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_FSEvents_002(app gopi.App, t *testing.T) {
	fsevents := app.UnitInstance("gopi/fsevents").(gopi.FSEvents)
	dir := os.TempDir()

	if fh, err := fsevents.Watch(dir, 0, func(watch uint32, evt gopi.FSEvent) {
		t.Log("watch=", watch, "evt=", evt)
	}); err != nil {
		t.Error(err)
	} else {
		// Wait for things to settle
		time.Sleep(time.Second)
		for i := 0; i < 100; i++ {
			// Create temporary folder and then remove it
			if tmpdir, err := ioutil.TempDir("", "fsevents"); err != nil {
				t.Error(err)
			} else if err := os.Remove(tmpdir); err != nil {
				t.Error(err)
			} else {
				t.Log("Created and removed", tmpdir)
			}
		}
		// Unwatch
		if err := fsevents.Unwatch(fh); err != nil {
			t.Error(err)
		}
	}
}
