/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi_test

import (
	"context"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi"
	logger "github.com/djthorpe/gopi/sys/logger"
	rpc "github.com/djthorpe/gopi/sys/rpc"
)

func TestRPCDiscovery_000(t *testing.T) {
	if logger, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Fatal(err)
	} else if driver, err := gopi.Open(rpc.Config{}, logger.(gopi.Logger)); err != nil {
		t.Fatal(err)
	} else {
		defer driver.Close()
		defer logger.Close()

		mdns := driver.(gopi.RPCServiceDiscovery)
		serviceType := "_smb._tcp"

		// Browse service records
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := mdns.Browse(ctx, serviceType); err != nil {
			t.Error(err)
		}

		// Register a service
		if err := mdns.Register(&gopi.RPCService{Name: "My Service", Type: serviceType, Port: 8000}); err != nil {
			t.Error(err)
		}

		// Wait for 5 seconds
		time.Sleep(5 * time.Second)

		// Cancel browsing
		cancel()
	}
}
