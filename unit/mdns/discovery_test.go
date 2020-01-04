/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns_test

import (
	"context"
	"fmt"
	"sync"
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

	app.Log().Debug("Calling Enumerate Services with 1s timeout")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if services, err := discovery.EnumerateServices(ctx); err != nil {
		return err
	} else {
		app.Log().Debug("Enumerate Services done, services=", services)
	}

	// Success
	return nil
}

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

	app.Log().Debug("Calling Enumerate Services with 1s timeout")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if names, err := discovery.EnumerateServices(ctx); err != nil {
		return err
	} else {
		var wait sync.WaitGroup
		for _, name := range names {
			wait.Add(1)
			go func(name string) {
				defer wait.Done()
				app.Log().Debug("Calling Lookup", name, "with 1s timeout")
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				if records, err := discovery.Lookup(ctx, name); err != nil {
					app.Log().Error(err)
				} else {
					for _, record := range records {
						fmt.Printf("%-20s %-60s %-20s:%d\n", record.Service, record.Name, record.Host, record.Port)
					}
				}
			}(name)
		}
		wait.Wait()
		app.Log().Debug("Lookup done")
	}

	// Success
	return nil
}
