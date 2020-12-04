package csv

import (
	"context"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Writer struct {
	sync.Mutex
	gopi.Unit
	gopi.Logger
	gopi.Publisher

	// Flags & Parameters
	path   *string
	ext    *string
	append *bool

	// Member variables
	files map[string]*file
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	runeComment = "#"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Writer) Define(cfg gopi.Config) error {
	this.path = cfg.FlagString("csv.path", "", "Metrics path")
	this.ext = cfg.FlagString("csv.ext", ".csv", "Metrics file extension")
	this.append = cfg.FlagBool("csv.append", true, "Append metrics to existing files")
	return nil
}

func (this *Writer) New(cfg gopi.Config) error {

	// Check path is a folder
	if *this.path != "" {
		if stat, err := os.Stat(*this.path); os.IsNotExist(err) {
			return gopi.ErrBadParameter.WithPrefix("-csv.path")
		} else if err != nil {
			return err
		} else if stat.IsDir() == false {
			return gopi.ErrBadParameter.WithPrefix("-csv.path")
		}
	}

	// Create file mapping
	this.files = make(map[string]*file)

	// Return success
	return nil
}

func (this *Writer) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Close all opened files
	for _, file := range this.files {
		file.Close()
	}

	// Release resources
	this.files = nil

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
			return ctx.Err()
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

	var result error

	// Create new files
	for _, metric := range metrics {
		key := metric.Name()
		if _, exists := this.files[key]; exists == false {
			if file, err := NewFile(*this.path, key, *this.ext, *this.append); err != nil {
				result = multierror.Append(result, err)
			} else {
				this.files[key] = file
			}
		}
	}

	// Write metrics
	for _, metric := range metrics {
		key := metric.Name()
		if err := this.files[key].Write(metric); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
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
