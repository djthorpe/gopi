/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package logger_test

import (
	"os"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	logger "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_Logger_000(t *testing.T) {
	t.Log("Test_Logger_000")
}

func Test_Logger_001(t *testing.T) {
	if logger, err := gopi.New(logger.Log{
		Writer: os.Stderr,
		Unit:   "logger_test",
	}, nil); err != nil {
		t.Error(err)
	} else {
		t.Log(logger)
	}
}

func Test_Logger_002(t *testing.T) {
	logger_, err := gopi.New(logger.Log{Writer: os.Stderr, Unit: "logger_test"}, nil)
	if err != nil {
		t.Error(err)
	}
	logger := logger_.(gopi.Logger)
	if logger.IsDebug() != false {
		t.Error("Expected IsDebug = false")
	}
	if logger.Name() != "logger_test" {
		t.Error("Expected Name = logger_test")
	}
}

func Test_Logger_003(t *testing.T) {
	logger_, err := gopi.New(logger.Log{Writer: os.Stderr, Unit: "logger_test", Debug: true}, nil)
	if err != nil {
		t.Error(err)
	}
	logger := logger_.(gopi.Logger)
	if logger.IsDebug() != true {
		t.Error("Expected IsDebug = false")
	}
}
