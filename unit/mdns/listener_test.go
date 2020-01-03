/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns_test

import (
	"testing"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
)

func Test_Listener_000(t *testing.T) {
	t.Log("Test_Listener_000")
}

func Test_Listener_001(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewDebugTool(Main_Test_Listener_001, flags, []string{"gopi/mdns/listener"}); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Listener_001(app gopi.App, _ []string) error {
	listener := app.UnitInstance("gopi/mdns/listener").(gopi.Publisher)
	if listener == nil {
		return gopi.ErrInternalAppError.WithPrefix("UnitInstance() failed")
	}
	app.Log().Debug(listener)
	time.Sleep(1 * time.Second)
	// Success
	return nil
}
