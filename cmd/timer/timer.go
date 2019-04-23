/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Runs either a one-shot or interval timer
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/timer"
)

type TimerType uint

const (
	TIMER_PERIODIC TimerType = iota
	TIMER_PERIODIC_IMMEDIATE
	TIMER_TIMEOUT
	TIMER_TIMEOUT2
	TIMER_BACKOFF
)

func (t TimerType) String() string {
	switch t {
	case TIMER_PERIODIC:
		return "TIMER_PERIODIC"
	case TIMER_PERIODIC_IMMEDIATE:
		return "TIMER_PERIODIC_IMMEDIATE"
	case TIMER_TIMEOUT:
		return "TIMER_TIMEOUT"
	case TIMER_TIMEOUT2:
		return "TIMER_TIMEOUT2"
	case TIMER_BACKOFF:
		return "TIMER_BACKOFF"
	default:
		return "[?? Invalid Value]"
	}
}

////////////////////////////////////////////////////////////////////////////////

func handleEvent(app *gopi.AppInstance, evt gopi.TimerEvent) {
	fmt.Println("EVENT: ", evt)

	// Schedule a new event when TIMER_TIMEOUT is fired
	if evt.UserInfo().(TimerType) == TIMER_TIMEOUT {
		app.Timer.NewTimeout(10*time.Second, TIMER_TIMEOUT2)
	}

	// Cancel all events which are over 20 fires
	if evt.Counter() > 20 {
		evt.Cancel()
	}
}

func Events(app *gopi.AppInstance, done <-chan struct{}) error {

	// Subscribe to timers
	edge := app.Timer.Subscribe()

FOR_LOOP:
	for {
		select {
		case evt := <-edge:
			if evt != nil {
				handleEvent(app, evt.(gopi.TimerEvent))
			}
		case <-done:
			break FOR_LOOP
		}
	}

	// Unsubscribe from events
	app.Timer.Unsubscribe(edge)
	return nil
}

func Main(app *gopi.AppInstance, done chan<- struct{}) error {

	// Schedule Interval timer which fires every second
	app.Timer.NewInterval(1*time.Second, TIMER_PERIODIC, false)

	// Schedule Interval timer which fires every other second, but which
	// also fires immediately
	app.Timer.NewInterval(2*time.Second, TIMER_PERIODIC_IMMEDIATE, true)

	// Schedule Timeout which fires once after 4 seconds
	app.Timer.NewTimeout(4*time.Second, TIMER_TIMEOUT)

	// Schedule Backoff timer which fires immediately then after 8 seconds,
	// and subsequently at backoff intervals up to 4 minutes
	app.Timer.NewBackoff(8*time.Second, 4*time.Minute, TIMER_BACKOFF)

	// wait until done
	app.Logger.Info("Waiting for CTRL+C")
	app.WaitForSignal()

	// Send done signal
	done <- gopi.DONE

	// Return nil
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the timer instance
	config := gopi.NewAppConfig("timer")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main, Events))
}
