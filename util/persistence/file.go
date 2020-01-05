/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

package persistence

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Config interface to return configuration values
type Config interface {
	// Returns the default filename for the persisent file
	DefaultFilename() string

	// Delay in writing to disk after modification
	WriteDelta() time.Duration

	// Path returns the path to the persistent file, or if empty
	// no writes happen to disk
	Path() string

	// Indent determines if the data should be written indented
	Indent() bool
}

// File implements filesystem persistence with JSON
type File struct {
	log              gopi.Logger
	filename_default string
	write_delta      time.Duration
	path             string
	modified         bool
	data             interface{}
	indent           bool

	sync.Mutex
	event.Tasks
}

////////////////////////////////////////////////////////////////////////////////
// INIT AND CLOSE

// Init initializes the persistence, passing in configuration, the data to
// persist and the logging module
func (this *File) Init(config Config, data interface{}, logger gopi.Logger) error {
	// Check incoming parameters
	if config == nil || logger == nil || data == nil {
		return gopi.ErrBadParameter
	}

	logger.Debug("<persistence.File>Init{ %+q }", config)

	// The data must be a pointer
	if typeof := reflect.TypeOf(data); typeof.Kind() != reflect.Ptr {
		return fmt.Errorf("Init: data must be a pointer to a data structure")
	}

	// Delta must be one or more seconds
	if delta := config.WriteDelta().Truncate(time.Second); delta > 0 {
		this.write_delta = delta
	} else {
		return gopi.ErrBadParameter
	}

	// Set member variables
	this.log = logger
	this.data = data
	this.indent = config.Indent()

	// Set default filename
	if filename_default := config.DefaultFilename(); filename_default == "" {
		return gopi.ErrBadParameter
	} else {
		this.filename_default = filename_default
	}

	// Read the persistence file
	if config.Path() != "" {
		if err := this.readPath(config.Path()); err != nil {
			return fmt.Errorf("Read: %v: %v", config.Path(), err)
		}
	}

	// Start process to write occasionally to disk
	this.Tasks.Start(this.writeTask)

	// Success
	return nil
}

// Close ends the persistence
func (this *File) Close() error {
	this.log.Debug("<persistence.File>Close{ path=%v }", strconv.Quote(this.path))

	// Stop all tasks
	if err := this.Tasks.Close(); err != nil {
		return err
	}

	// Success
	return nil
}

// String stringifies
func (this *File) String() string {
	return fmt.Sprintf("<persistence.File>{ path=%v modified=%v indent=%v }", strconv.Quote(this.path), this.isModified(), this.indent)
}

////////////////////////////////////////////////////////////////////////////////
// READ

// ReadPath creates regular file if it doesn't exist, or else reads from the path
func (this *File) readPath(path string) error {
	this.log.Debug2("<persistence.File>ReadPath{ path=%v }", strconv.Quote(path))

	// Append home directory if relative path
	if filepath.IsAbs(path) == false {
		if homedir, err := os.UserHomeDir(); err != nil {
			return err
		} else {
			path = filepath.Join(homedir, path)
		}
	}

	// Set path
	this.path = path

	// Append filename
	if stat, err := os.Stat(this.path); err == nil && stat.IsDir() {
		// append default filename
		this.path = filepath.Join(this.path, this.filename_default)
	}

	// Read file
	if stat, err := os.Stat(this.path); err == nil && stat.Mode().IsRegular() {
		if err := this.readPath_(this.path); err != nil {
			return err
		} else {
			return nil
		}
	} else if os.IsNotExist(err) {
		// Create file
		if fh, err := os.Create(this.path); err != nil {
			return err
		} else if err := fh.Close(); err != nil {
			return err
		} else {
			this.SetModified()
			return nil
		}
	} else {
		return err
	}
}

func (this *File) readPath_(path string) error {
	this.Lock()
	defer this.Unlock()

	// Check for zero-sized file, and don't read if it is zero-sized
	if stat, err := os.Stat(path); err != nil {
		return err
	} else if stat.Size() == 0 {
		return nil
	}

	// Open and read the JSON
	if fh, err := os.Open(path); err != nil {
		return err
	} else {
		defer fh.Close()
		if err := this.Reader(fh); err != nil {
			return err
		} else {
			this.modified = false
		}
	}

	// Success
	return nil
}

// Reader reads the configuration from an io.Reader object
func (this *File) Reader(fh io.Reader) error {
	dec := json.NewDecoder(fh)
	if err := dec.Decode(&this.data); err != nil {
		return err
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// WRITE

// writePath writes the configuration file to disk
func (this *File) writePath(path string) error {
	this.log.Debug2("<persistence.File>WritePath{ path=%v }", strconv.Quote(path))
	this.Lock()
	defer this.Unlock()
	if fh, err := os.Create(path); err != nil {
		return err
	} else {
		defer fh.Close()
		if err := this.Writer(fh); err != nil {
			return err
		} else {
			this.modified = false
		}
	}

	// Success
	return nil
}

// Writer writes to io.Writer object
func (this *File) Writer(fh io.Writer) error {
	enc := json.NewEncoder(fh)
	if this.indent {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(this.data); err != nil {
		return err
	}
	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SET MODIFIED FLAG

// SetModified sets the modified flag to true
func (this *File) SetModified() {
	this.log.Debug2("<persistence.File>SetModified{}")
	this.Lock()
	defer this.Unlock()
	this.modified = true
}

// isModified returns the modified flag
func (this *File) isModified() bool {
	return this.modified
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND TASKS

func (this *File) writeTask(start chan<- event.Signal, stop <-chan event.Signal) error {
	start <- gopi.DONE
	ticker := time.NewTimer(100 * time.Millisecond)
FOR_LOOP:
	for {
		select {
		case <-ticker.C:
			if this.isModified() {
				if this.path == "" {
					// Do nothing
				} else if err := this.writePath(this.path); err != nil {
					this.log.Warn("Write: %v: %v", this.path, err)
				}
			}
			ticker.Reset(this.write_delta)
		case <-stop:
			break FOR_LOOP
		}
	}

	// Stop the ticker
	ticker.Stop()

	// Try and write
	if this.modified {
		if this.path == "" {
			// Do nothing
		} else if err := this.writePath(this.path); err != nil {
			this.log.Warn("Write: %v: %v", this.path, err)
		}
	}

	// Success
	return nil
}
