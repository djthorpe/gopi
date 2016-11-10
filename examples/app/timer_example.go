/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Define an interface which generates events from an external source
type EventGenerator interface {
	// Start runloop for the generator
	Run(chan Event, chan bool, chan bool)
}

// This is the event structure
type Event struct {
	Sender    EventGenerator
	Timestamp time.Time
}

// Timer is a concrete EventGenerator
type Timer struct {
	After       time.Duration
	FireAtStart bool
	Repeating   bool
}

type App struct {
	signal_channel chan os.Signal
	done_channel   chan bool
	finish_channel chan bool
	event_channel  chan Event
}

type SignalCallback func()

type EventCallback func(*Event)

// Implement EventGenerator
func (this *Timer) Run(event_channel chan Event, finish_channel chan bool, done_channel chan bool, callback EventCallback) {
	done := false
	if this.FireAtStart {
		// Fire the event
		event_channel <- evt{Event{Sender: this, Callback}, callback}
	}
	for done == false {
		select {
		case <-finish_channel:
			// If we receive a finish signal, then break
			done = true
			break
		case <-time.After(this.After):
			// Fire the event
			event_channel <- Event{Sender: this}
			if this.Repeating == false {
				done = true
			}
		}
	}
	done_channel <- done
}

func NewApp() (*App, error) {
	this := new(App)
	this.signal_channel = make(chan os.Signal, 1)
	this.finish_channel = make(chan bool, 1)
	this.event_channel = make(chan Event, 1)
	this.done_channel = make(chan bool, 1)

	// Go routine to wait for signal, and send finish signal in that case
	go func() {
		<-this.signal_channel
		this.finish_channel <- true
	}()

	// Success
	return this, nil
}

func (this *App) AddEventGenerator(generator EventGenerator, callback EventCallback) {
	// start event generator
	go generator.Run(this.event_channel, this.finish_channel, this.done_channel, callback)
}

func (this *App) ScheduleTimerWithInterval(interval time.Duration, repeats bool, callback EventCallback) {
	this.AddEventGenerator(&Timer{After: interval, Repeating: repeats}, callback)
}

func (this *App) CatchSignals(callback SignalCallback, signals ...os.Signal) {
	signal.Stop(this.signal_channel)
	signal.Notify(this.signal_channel, signals...)
}

func (this *App) Run() {
	// Runloop accepting events
	done := false
	for done == false {
		select {
		case event := <-this.event_channel:
			fmt.Println("PING", event)
		case done = <-this.done_channel:
			fmt.Println("DONE", done)
			break
		}
	}
}

func (this *App) Stop() {
	fmt.Println("STOP APP")
}

// Main function
func main() {
	app, err := NewApp()
	if err != nil {
		fmt.Println(err)
		return
	}

	// add repeating timer
	app.ScheduleTimerWithInterval(time.Second*5.0, true, func(event *Event) {
		fmt.Println("TIMER 1 FIRED, event = ", event)
	})

	// catch signal - and then call application to stop
	app.CatchSignals(func() {
		app.Stop()
	}, syscall.SIGTERM, syscall.SIGINT)

	// run the app - should return when a signal is caught and app.Stop() is run
	fmt.Println("START RUN")
	app.Run()
	fmt.Println("STOP RUN")
}
