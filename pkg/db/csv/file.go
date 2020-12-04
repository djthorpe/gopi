package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type file struct {
	sync.Mutex

	path   string
	fh     *os.File
	writer *csv.Writer
}

////////////////////////////////////////////////////////////////////////////////
// NEW AND CLOSE

func NewFile(path, name, ext string, append bool) (*file, error) {
	this := new(file)
	name = name + "." + strings.Trim(ext, ".")

	// Set the path for writing
	if path, err := filepath.Abs(filepath.Join(path, name)); err != nil {
		return nil, err
	} else {
		this.path = path
	}

	// Create or append the file
	mode := os.O_WRONLY | os.O_CREATE
	if append {
		mode |= os.O_APPEND
	}
	if fh, err := os.OpenFile(this.path, mode, 0666); err != nil {
		return nil, err
	} else {
		this.fh = fh
	}

	// Create the writer
	this.writer = csv.NewWriter(this.fh)

	// Return success
	return this, nil
}

func (this *file) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var result error

	// Flush any existing data
	this.writer.Flush()
	if err := this.writer.Error(); err != nil {
		result = multierror.Append(result, err)
	}

	// Close filehandle
	if this.fh != nil {
		if err := this.fh.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.fh = nil
	this.writer = nil

	// Return errors
	return result
}

func (this *file) Write(metric gopi.Measurement) error {
	// Check size of file, so that empty files
	// can have header and comment written
	size := int64(0)
	if stat, err := this.fh.Stat(); err == nil {
		size = stat.Size()
	}

	// Generate the header, comment and row
	header := []string{}
	comment := []string{}
	row := []string{}
	if t := metric.Time(); t.IsZero() == false {
		header = append(header, "time")
		comment = append(comment, "time")
		row = append(row, t.Format(time.RFC3339))
	}
	for _, tag := range metric.Tags() {
		header = append(header, tag.Name())
		comment = append(comment, "tag["+tag.Kind()+"]")
		if tag.IsNil() {
			row = append(row, "")
		} else {
			row = append(row, fmt.Sprint(tag.Value()))
		}
	}
	for _, metric := range metric.Metrics() {
		header = append(header, metric.Name())
		comment = append(comment, "metric["+metric.Kind()+"]")
		if metric.IsNil() {
			row = append(row, "")
		} else {
			row = append(row, fmt.Sprint(metric.Value()))
		}
	}

	// Write the header
	if size == 0 {
		if err := this.writer.Write(header); err != nil {
			return err
		}
		comment[0] = string(runeComment) + " " + comment[0]
		if err := this.writer.Write(comment); err != nil {
			return err
		}
	}

	// Write the row
	if err := this.writer.Write(row); err != nil {
		return err
	}

	// Flush and return any errors
	this.writer.Flush()
	return this.writer.Error()
}
