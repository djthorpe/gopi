package tool

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/config"
	"github.com/djthorpe/gopi/v3/pkg/graph"
)

func CommandLine(name string, args []string, objs ...interface{}) int {
	// Create empty configuration and graph
	cfg := config.New(name, args)
	graph, err := graph.Create(objs...)
	if err != nil {
		fmt.Fprintln(os.Stderr, "New:", err)
		return -1
	}

	// Get logger object
	logger := graph.GetLogger()

	// Call Define for each object
	if err := graph.Define(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "Define:", err)
		return -1
	}

	// Parse command-line arguments
	if err := cfg.Parse(); errors.Is(err, gopi.ErrHelp) || errors.Is(err, flag.ErrHelp) {
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Config:", err)
		return -1
	}

	// Call New
	if err := graph.New(cfg); errors.Is(err, gopi.ErrHelp) || errors.Is(err, flag.ErrHelp) {
		cfg.Usage("")
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "New:", err)
		return -1
	}

	// If there is a gopi.Logger object and debug is set then
	// use the Debug method to output extra information
	if logger != nil && logger.IsDebug() {
		graph.Logfn = logger.Debug
	}

	// Call Dispose on exit
	defer func() {
		if err := graph.Dispose(); err != nil {
			fmt.Fprintln(os.Stderr, "Dispose:", err)
		}
	}()

	// Create context with a cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Handle signals - call cancel when interrupt received
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		s := <-ch
		if logger != nil && logger.IsDebug() {
			logger.Debug("Got signal: ", s)
		}
		cancel()
	}()

	// Call Run
	if err := graph.Run(ctx); err != nil && err != context.Canceled {
		fmt.Fprintln(os.Stderr, "Run:", err)
		return -1
	}

	// Return success
	return 0
}
