package tool

import (
	"context"
	"errors"
	"flag"
	"reflect"
	"sync"
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/config"
	"github.com/djthorpe/gopi/v3/pkg/graph"
)

func Test(t *testing.T, args []string, obj, fn interface{}) int {
	// Create empty configuration and graph
	cfg := config.New(t.Name(), args)
	g := graph.NewGraph(t.Log)

	// Create objects
	if err := g.Create(obj); err != nil {
		t.Error("New:", err)
		return -1
	}

	// Call Define for each object
	if err := g.Define(cfg); err != nil {
		t.Error("Define:", err)
		return -1
	}

	// Parse command-line arguments
	if err := cfg.Parse(); errors.Is(err, gopi.ErrHelp) || errors.Is(err, flag.ErrHelp) {
		return 0
	} else if err != nil {
		t.Error("Config:", err)
		return -1
	}

	// Call New
	if err := g.New(cfg); errors.Is(err, gopi.ErrHelp) || errors.Is(err, flag.ErrHelp) {
		cfg.Usage("")
		return 0
	} else if err != nil {
		t.Error("New:", err)
		return -1
	}

	var wg sync.WaitGroup

	// Create context with a cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Run tests here, and call cancel when done to end the 'Run'
	// functions
	if fn_ := reflect.ValueOf(fn); fn_.Kind() != reflect.Func {
		t.Error("Invalid test function")
		return -1
	} else {
		wg.Add(1)
		go func() {
			fn_.Call([]reflect.Value{reflect.ValueOf(obj)})
			cancel()
			wg.Done()
		}()
	}

	// Call run
	if err := g.Run(ctx); err != nil && err != context.Canceled {
		t.Error("Run:", err)
		return -1
	}

	// Wait for end of test routine before dispose called
	wg.Wait()

	// Call Dispose
	if err := g.Dispose(); err != nil {
		t.Error("Dispose:", err)
		return -1
	}

	// Return success
	return 0
}
