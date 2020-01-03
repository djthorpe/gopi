/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns_test

import (
	"context"
	"net"
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

func Test_Register_000(t *testing.T) {
	t.Log("Test_Register_000")
}

func Test_Register_001(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewDebugTool(Main_Test_Register_001, flags, []string{"register"}); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Register_001(app gopi.App, _ []string) error {
	register := app.UnitInstance("register").(gopi.RPCServiceRegister)
	if register == nil {
		return gopi.ErrInternalAppError.WithPrefix("UnitInstance() failed")
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		addr := net.ParseIP("127.0.0.1")
		if err := register.Register(ctx, gopi.RPCServiceRecord{
			Name:    "First Test",
			Service: "_gopi._tcp",
			Host:    "test1",
			Port:    8080,
			Txt:     []string{"name=test1"},
			Addrs:   []net.IP{addr},
		}); err != nil {
			app.Log().Error(err)
		}
	}()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		addr := net.ParseIP("127.0.0.1")
		if err := register.Register(ctx, gopi.RPCServiceRecord{
			Name:    "Second Test",
			Service: "_gopi._tcp",
			Host:    "test2",
			Port:    8080,
			Txt:     []string{"name=test2"},
			Addrs:   []net.IP{addr},
		}); err != nil {
			app.Log().Error(err)
		}
	}()

	time.Sleep(30 * time.Second)

	// Success
	return nil
}
