package mdns_test

import (
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
	*mdns.Discovery
}

func Test_Discovery_001(t *testing.T) {
	tool.Test(t, nil, new(DiscoveryApp), func(app *DiscoveryApp) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := app.Discovery.EnumerateServices(ctx); err != nil {
			t.Error(err)
		}

	})
}
