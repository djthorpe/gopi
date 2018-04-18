
## Events and Timers

The framework relies heavily on being __event-driven__ in many of
the drivers. For example,

  * The `gpio` module will fire events when an input is changed from
    high to low (falling edge) or low to high (rising edge);
  * The `input` module will fire events on key press, mouse move or
    touchscreen press;
  * The `lirc` module will fire events when infrared light is detected;
  * The `timer` module can fire events once or repeatedly at set intervals.

Your code can subscribe to events and react accordingly. There are two
interfaces which are implemented:

  * The `gopi.Publisher` interface implements subscribe, unsubscribe and
    emit methods;
  * The `gopi.Event` iterface implements an event which can be emitted.

In addition there are two utility structures:

  * The `event.PubSub` structure implements a "fan out" for emitting events
    to one or more subscriber channels;
  * The 'event.Merge` structure implements a "fan in" for merging events
    from several channels into a single channel.

Here is an example of subscribing to events from two modules, and reacting
to them:

```
func eventLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	// Subscribe to timers and GPIO edges
	timer_chan := app.Timer.Subscribe()
    gpio_chan := app.GPIO.Subscribe()

FOR_LOOP:
	for {
		select {
		case evt := <-timer_chan:
			handleTimerEvent(app, evt.(gopi.TimerEvent))
        case evt := <-gpio_chan:
			handleGPIOEvent(app, evt.(gopi.GPIOEvent))
		case <-done:
			break FOR_LOOP
		}
	}

	// Unsubscribe from events
	app.Timer.Unsubscribe(timer_chan)
	app.GPIO.Unsubscribe(gpio_chan)

	return nil
}
```



