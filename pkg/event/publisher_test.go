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

		var wg, wg2 sync.WaitGroup

		// Number of messages to send
		msgs := rand.Intn(1000) + 1

		// Subscribe
		ch := app.Publisher.Subscribe()
		defer app.Publisher.Unsubscribe(ch)

		// Emit all null events
		wg2.Add(1)
		go func(n int) {
			wg.Add(n)
			for i := 0; i < n; i++ {
				t.Log("Emitting events", i+1, "of", n)
				if err := app.Publisher.Emit(nil, true); err != nil {
					t.Error(err)
				}
			}
			wg2.Done()
		}(msgs)

		// Receive null events
		wg2.Add(1)
		go func(n int) {
			i := 1
			for _ = range ch {
				t.Log("Receiving event", i, "of", n)
				i++
				wg.Done()
			}
			wg2.Done()
		}(msgs)

		t.Log("Waiting for all messages received")
		wg.Wait()
		t.Log("Waiting for goroutines to end")
		//wg2.Wait()

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
