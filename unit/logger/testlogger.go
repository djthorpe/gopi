/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package logger

import (
	"os"
	"sync"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type TestLogger struct {
	Unit string
	T    *testing.T
}

type testlogger struct {
	T *testing.T

	base.Logger
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (TestLogger) Name() string { return "gopi/testlogger" }

func (config TestLogger) New(gopi.Logger) (gopi.Unit, error) {
	this := new(testlogger)
	if err := this.Logger.Init(os.Stderr, config.Unit, true); err != nil {
		return nil, err
	} else if config.T == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("T")
	} else {
		this.T = config.T
	}

	// Return success
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Logger

func (this *testlogger) Error(err error) error {
	this.Lock()
	defer this.Unlock()
	this.T.Error(err)
	return err
}

func (this *testlogger) Warn(args ...interface{}) {
	this.Lock()
	defer this.Unlock()
	this.T.Log(append([]interface{}{"WARN:"}, args...)...)
}

func (this *testlogger) Info(args ...interface{}) {
	this.Lock()
	defer this.Unlock()
	this.T.Log(append([]interface{}{"INFO:"}, args...)...)
}

func (this *testlogger) Debug(args ...interface{}) {
	this.Lock()
	defer this.Unlock()
	this.T.Log(append([]interface{}{"DEBUG:"}, args...)...)
}
