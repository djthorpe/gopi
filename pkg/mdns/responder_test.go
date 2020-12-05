package mdns_test

import (
	"context"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	mdns "github.com/djthorpe/gopi/v3/pkg/mdns"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/event"
)

type ResponderApp struct {
	gopi.Unit
	*mdns.Responder
}

func (this *ResponderApp) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func Test_Responder_001(t *testing.T) {
	tool.Test(t, nil, new(ResponderApp), func(app *ResponderApp) {
		if app.Responder == nil {
			t.Error("Expected non-nil Responder")
		}
	})
}

func Test_Responder_002(t *testing.T) {
	tool.Test(t, nil, new(ResponderApp), func(app *ResponderApp) {
		good := []struct{ s, n string }{
			{"_test._tcp.", "Service Name"},
			{"_test._tcp.", "Special with é accents"},
		}
		for _, test := range good {
			r, err := app.Responder.NewServiceRecord(test.s, test.n, 1, nil, gopi.SERVICE_FLAG_IP4)
			if err != nil {
				t.Error(err)
			} else if r == nil {
				t.Error("Unexpected return from NewServiceRecord")
			} else if s2 := r.Service(); s2 != test.s {
				t.Error("Unexpected service", s2, "(expected", test.s, ")")
			} else if n2 := r.Name(); n2 != test.n {
				t.Error("Unexpected name", n2, "(expected", test.n, ")")
			} else {
				t.Log(test, "=>", r)
			}
		}
	})
}

func Test_Responder_003(t *testing.T) {
	tool.Test(t, nil, new(ResponderApp), func(app *ResponderApp) {
		r, err := app.Responder.NewServiceRecord("_gopi._tcp", t.Name(), 9999, nil, gopi.SERVICE_FLAG_IP4)
		if err != nil {
			t.Error(err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := app.Responder.Serve(ctx, []gopi.ServiceRecord{r}); err != nil {
			t.Error(err)
		}
	})
}
