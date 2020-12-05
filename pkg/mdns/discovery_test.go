package mdns_test

import (
	"context"
	"sync"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/event"
	_ "github.com/djthorpe/gopi/v3/pkg/mdns"
)

type DiscoveryApp struct {
	gopi.Unit
	gopi.Logger
	gopi.ServiceDiscovery
}

func (this *DiscoveryApp) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func Test_Discovery_001(t *testing.T) {
	tool.Test(t, nil, new(DiscoveryApp), func(app *DiscoveryApp) {
		if app.ServiceDiscovery == nil {
			t.Error("No ServiceDiscovery object")
		}
	})
}

func Test_Discovery_002(t *testing.T) {
	tool.Test(t, nil, new(DiscoveryApp), func(app *DiscoveryApp) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Cancel after one second
		if services, err := app.ServiceDiscovery.EnumerateServices(ctx); err != nil {
			t.Error("EnumerateServices:", err)
		} else {
			t.Log("EnumerateServices:", services)
		}
	})
}

func Test_Discovery_003(t *testing.T) {
	tool.Test(t, nil, new(DiscoveryApp), func(app *DiscoveryApp) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Cancel after one second
		if services, err := app.ServiceDiscovery.EnumerateServices(ctx); err != nil {
			t.Error("EnumerateServices:", err)
		} else {
			var wg sync.WaitGroup
			for _, service := range services {
				wg.Add(1)
				go func(service string) {
					defer wg.Done()
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					defer cancel()
					if r, err := app.ServiceDiscovery.Lookup(ctx, service); err != nil {
						t.Error(err)
					} else {
						t.Log("Lookup:", service, r)
					}
				}(service)
			}
			wg.Wait()
		}
	})
}
