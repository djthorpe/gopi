/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bus_test

import (
	"context"
	"testing"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_Bus_000(t *testing.T) {
	t.Log("Test_Bus_000")
}

func Test_Bus_001(t *testing.T) {
	gotEvent := false
	main := func(app gopi.App, t *testing.T) {
		t.Log(app)
		bus := app.Bus()
		if bus == nil {
			t.Fatal(gopi.ErrInternalAppError.WithPrefix("Missing Bus()"))
		}
		t.Log("-> RUN()", app.Bus())

		// Set a default handler
		bus.DefaultHandler(gopi.EVENT_NS_DEFAULT, func(_ context.Context, _ gopi.App, evt gopi.Event) {
			t.Log("-> EVENT()", evt)
			if evt == gopi.NullEvent {
				gotEvent = true
			}
			// Simulate event taking a while to handle...
			time.Sleep(time.Second)
			t.Log("<- EVENT()")
		})

		// Emit null event
		bus.Emit(gopi.NullEvent)
		// End of run
		t.Log("<- RUN()")
	}
	args := []string{"-debug"}
	units := []string{"bus"}

	if app, err := app.NewTestTool(t, main, args, units...); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	} else if gotEvent == false {
		t.Error("Emitted event was not handled")
	}
}

func Test_Bus_002(t *testing.T) {
	var gotEvent = false
	if app, err := app.NewCommandLineTool(func(app gopi.App, _ []string) error {
		bus := app.Bus()
		if bus == nil {
			return gopi.ErrInternalAppError.WithPrefix("Missing Bus()")
		}
		t.Log("-> RUN()", app.Bus())

		// Set a default handler
		bus.NewHandler(gopi.EventHandler{Name: "gopi.NullEvent", Handler: func(_ context.Context, _ gopi.App, evt gopi.Event) {
			t.Log("-> EVENT()", evt)
			if evt == gopi.NullEvent {
				gotEvent = true
			}
			// Simulate event taking a while to handle...
			time.Sleep(time.Second)
			t.Log("<- EVENT()")
		}})

		// Emit null event
		bus.Emit(nil)

		// End of run
		t.Log("<- RUN()")

		return nil
	}, nil, "bus"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	} else if gotEvent == false {
		t.Error("Emitted event was not handled")
	}
}

func Test_Bus_003(t *testing.T) {
	var gotEvent = false
	var gotTimeout = false
	if app, err := app.NewCommandLineTool(func(app gopi.App, _ []string) error {
		bus := app.Bus()
		timeout := 2 * time.Second
		if bus == nil {
			return gopi.ErrInternalAppError.WithPrefix("Missing Bus()")
		}
		t.Log("-> RUN()", app.Bus())

		// Set a default handler
		bus.NewHandler(gopi.EventHandler{"gopi.NullEvent", func(ctx context.Context, _ gopi.App, evt gopi.Event) {
			t.Log("-> EVENT()", evt)
			then := time.Now()
			if evt == gopi.NullEvent {
				gotEvent = true
			}
			// Simulate event taking a while to handle (longer than timeout)
			timer := time.NewTimer(timeout * 2)
			select {
			case <-ctx.Done():
				delta := time.Now().Sub(then).Truncate(time.Second)
				t.Log("GOT CTX DONE AFTER", delta)
				gotTimeout = true
				if delta != timeout {
					t.Error("Unexpected timeout delta", delta, "expected", timeout)
				}
			case <-timer.C:
				t.Log("GOT TIMER DONE AFTER", time.Now().Sub(then).Truncate(time.Second))
			}
			t.Log("<- EVENT()")
		}, gopi.EVENT_NS_DEFAULT, timeout})

		// Emit null event
		bus.Emit(nil)

		// End of run
		t.Log("<- RUN()")

		return nil
	}, nil, "bus"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	} else if gotEvent == false {
		t.Error("Emitted event was not handled")
	} else if gotTimeout == false {
		t.Error("Emitted event did not timeout")
	}
}
