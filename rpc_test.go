package gopi_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi"
	logger "github.com/djthorpe/gopi/sys/logger"
	mdns "github.com/djthorpe/gopi/sys/mdns"
)

func TestRPCDiscovery_000(t *testing.T) {
	if logger, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Fatal(err)
	} else if driver, err := gopi.Open(mdns.Config{}, logger.(gopi.Logger)); err != nil {
		t.Fatal(err)
	} else {
		defer driver.Close()

		mdns := driver.(gopi.RPCDiscovery)
		serviceType := "_workstation._tcp"

		// Wait for service records
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		if err := mdns.Browse(ctx, serviceType); err != nil {
			t.Error(err)
		} else {
			logger.(gopi.Logger).Info("Browse completed")
		}

		// Register a service
		if err := mdns.Register("My Service", serviceType, 8000, nil); err != nil {
			t.Error(err)
		}

		fmt.Println(mdns)

		time.Sleep(20 * time.Second)
	}
}
