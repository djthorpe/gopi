/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns_test

import (
	"context"
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

func Test_Discovery_000(t *testing.T) {
	t.Log("Test_Discovery_000")
}

func Test_Discovery_001(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewDebugTool(Main_Test_Discovery_001, flags, []string{"discovery"}); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Discovery_001(app gopi.App, _ []string) error {
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	if discovery == nil {
		return gopi.ErrInternalAppError.WithPrefix("UnitInstance() failed")
	}
	app.Log().Debug(discovery)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if services, err := discovery.EnumerateServices(ctx); err != nil {
		return err
	} else {
		app.Log().Debug("services=", services)
	}

	// Success
	return nil
}

/*
func Test_Discovery_002(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewDebugTool(Main_Test_Discovery_002, flags, []string{"discovery"}); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Discovery_002(app gopi.App, _ []string) error {
	discovery := app.UnitInstance("discovery").(gopi.RPCServiceDiscovery)
	if discovery == nil {
		return gopi.ErrInternalAppError.WithPrefix("UnitInstance() failed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if services, err := discovery.EnumerateServices(ctx); err != nil {
		return err
	} else {
		app.Log().Debug("services=", services)
	}

	// Success
	return nil
}
*/
