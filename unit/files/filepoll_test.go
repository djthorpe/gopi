/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files_test

import (
	"encoding/binary"
	"fmt"
	"os"
	"testing"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/files"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

var (
	fdmap = make(map[uintptr]*os.File)
)

func Test_Filepoll_000(t *testing.T) {
	t.Log("Test_Filepoll_000")
}

func Test_Filepoll_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Filepoll_001, []string{"-debug"}, "gopi/filepoll"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Filepoll_001(app gopi.App, t *testing.T) {
	filepoll := app.UnitInstance("gopi/filepoll").(gopi.FilePoll)
	t.Log(filepoll)
}

func Test_Filepoll_002(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Filepoll_002, []string{"-debug"}, "gopi/filepoll"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Filepoll_002(app gopi.App, t *testing.T) {
	filepoll := app.UnitInstance("gopi/filepoll").(gopi.FilePoll)
	if tmp, err := os.Open("/dev/lirc1"); os.IsNotExist(err) {
		t.Log("Not running test, missing device")
	} else if err != nil {
		t.Error(err)
	} else {
		fdmap[tmp.Fd()] = tmp
		if err := filepoll.Watch(tmp.Fd(), gopi.FILEPOLL_FLAG_READ, Watcher_Test_Filepoll_002); err != nil {
			t.Error(err)
		} else {
			fmt.Println("Watching for events on", tmp.Name())
			time.Sleep(time.Second * 5)
			fmt.Println("End watching for events on", tmp.Name())
			if err := filepoll.Unwatch(tmp.Fd()); err != nil {
				t.Error(err)
			}
			fmt.Println("Unwatched", tmp.Name())
		}
	}
}

func Watcher_Test_Filepoll_002(fd uintptr, flags gopi.FilePollFlags) {
	handle := fdmap[fd]
	switch flags {
	case gopi.FILEPOLL_FLAG_READ:
		var value uint32
		if err := binary.Read(handle, binary.LittleEndian, &value); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(value)
		}
	}
}
