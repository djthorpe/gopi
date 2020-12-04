package mdns_test

import (
	"sync"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/event"
	mdns "github.com/djthorpe/gopi/v3/pkg/mdns"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
	"golang.org/x/net/context"
)

type DiscoveryApp struct {
	gopi.Unit
	gopi.Logger
	*mdns.Discovery
}

func (this *DiscoveryApp) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func Test_Discovery_001(t *testing.T) {
	tool.Test(t, nil, new(DiscoveryApp), func(app *DiscoveryApp) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Cancel after one second
		if services, err := app.Discovery.EnumerateServices(ctx); err != nil {
			t.Error("EnumerateServices:", err)
		} else {
			t.Log("EnumerateServices:", services)
		}
	})
}

func Test_Discovery_002(t *testing.T) {
	tool.Test(t, nil, new(DiscoveryApp), func(app *DiscoveryApp) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		// Cancel after one second
		if services, err := app.Discovery.EnumerateServices(ctx); err != nil {
			t.Error("EnumerateServices:", err)
		} else {
			var wg sync.WaitGroup
			for _, service := range services {
				wg.Add(1)
				go func(service string) {
					defer wg.Done()
					ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
					defer cancel()
					if r, err := app.Discovery.Lookup(ctx, service); err != nil {
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
