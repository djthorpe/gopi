package persistence_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util/persistence"

	// Modules
	logger "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

type tester struct {
	string_data string
	persistence.File
}

type tester_config struct {
}

func (tester_config) DefaultFilename() string {
	return "tester_config.json"
}

func (tester_config) WriteDelta() time.Duration {
	return time.Second * 5
}

func (tester_config) Path() string {
	return os.TempDir()
}

func (tester_config) Indent() bool {
	return true
}

func NewTester() (*tester, error) {
	this := new(tester)
	config := new(tester_config)
	if log, err := gopi.Open(logger.Config{}, nil); err != nil {
		return nil, err
	} else if err := this.File.Init(config, &this.string_data, log.(gopi.Logger)); err != nil {
		return nil, err
	} else {
		return this, nil
	}
}

func (this *tester) Close() error {
	return this.File.Close()
}

func (this *tester) String() string {
	return fmt.Sprintf("<tester>{ %v }", this.File.String())
}

////////////////////////////////////////////////////////////////////////////////

func TestFile_000(t *testing.T) {
	if store, err := NewTester(); err != nil {
		t.Fatal(err)
	} else {
		defer store.Close()
	}
}

func TestFile_001(t *testing.T) {
	if store, err := NewTester(); err != nil {
		t.Fatal(err)
	} else {
		defer store.Close()
		store.string_data = "Hello, World"
		store.SetModified()
	}
}
