package csv

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Writer struct {
	sync.Mutex
	gopi.Unit
	gopi.Logger
	gopi.Publisher

	// Flags & Parameters
	path *string
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Writer) Define(cfg gopi.Config) error {
	this.path = cfg.FlagString("csv.path", "", "Metrics Folder")
	return nil
}

func (this *Writer) New(cfg gopi.Config) error {

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *Writer) Run(ctx context.Context) error {
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	for {
		select {
		case evt := <-ch:
			if m, ok := evt.(gopi.Measurement); ok {
				if err := this.Write(m); err != nil {
					this.Print(err)
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Writer) Ping() (time.Duration, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// CSV writer does not have ping
	return 0, nil
}

// Write measurements to the endpoint
func (this *Writer) Write(metrics ...gopi.Measurement) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Return bad parameter if no metrics
	if len(metrics) == 0 {
		return gopi.ErrBadParameter.WithPrefix("Write")
	}

	// Noop
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Writer) String() string {
	str := "<writer.csv"
	if *this.path != "" {
		str += " path=" + strconv.Quote(*this.path)
	}
	return str + ">"
}
