/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package timer_test

import (
	"os"
	"testing"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"

	// Units
	logger "github.com/djthorpe/gopi/v2/unit/logger"
	timer "github.com/djthorpe/gopi/v2/unit/timer"
)

func Test_Timer_000(t *testing.T) {
	t.Log("Test_Timer_000")
}

func Test_Timer_001(t *testing.T) {
	config := timer.Timer{}
	// Create logger
	logger_, err := gopi.New(logger.Log{
		Writer: os.Stderr,
		Unit:   config.Name(),
		Debug:  true,
	}, nil)
	if err != nil {
		t.Error(err)
	}
	// Create timer
	timer, err := gopi.New(config, logger_.(gopi.Logger))
	if err != nil {
		t.Error(err)
	}
	defer timer.Close()
	// Create a ticker
	if tickerId := timer.(gopi.Timer).NewTicker(time.Second); tickerId == 0 {
		t.Error("NewTicker returned zero")
	} else {
		t.Log("timerId = ", tickerId)
		time.Sleep(10 * time.Second)
	}
}
