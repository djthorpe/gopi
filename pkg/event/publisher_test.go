package event

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

type App struct {
	gopi.Unit
	gopi.Publisher
}

func Test_Event_000(t *testing.T) {
	pub := &publisher{}
	if ch := pub.Subscribe(); ch == nil {
		t.Error("Unexpected nil return value")
	} else {
		pub.Unsubscribe(ch)
	}
}

func Test_Event_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Publisher == nil {
			t.Error("Unexpected nil value for publisher")
		}

		var wg sync.WaitGroup
		wg.Add(1)

		// Receive null events
		go func() {
			i := 0
			ch := app.Subscribe()
			wg.Done()
			for evt := range ch {
				t.Log("->Receive", i, evt)
				wg.Done()
				i++
			}
		}()

		wg.Wait()

		// Emit null events
		for i := 0; i < rand.Intn(1000)+1; i++ {
			wg.Add(1)
			t.Log("->Emit", i)
			app.Emit(nil, true)
		}
		wg.Wait()
	})
}

func Test_Event_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		var wg sync.WaitGroup
		wg.Add(1)

		// Receive null events
		go func() {
			i := 0
			ch := app.Subscribe()
			wg.Done()
			for evt := range ch {
				t.Log("->Receive", i, evt)
				wg.Done()
				i++
			}
		}()

		wg.Wait()

		// Emit null events
		for i := 0; i < rand.Intn(1000)+1; i++ {
			wg.Add(1)
			for {
				t.Log("->Emit", i)
				if err := app.Emit(nil, false); err != nil {
					t.Log("  (channel full, waiting for 10ms)")
					time.Sleep(time.Millisecond * 10)
				} else {
					break
				}
			}
		}
		wg.Wait()
	})
}
