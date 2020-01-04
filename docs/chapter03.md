# Event Handling

Fundamentally the `gopi` framework implements tools by reacting to events. The `gopi.Event` interface defines a basic event to be handled:

```go
type Event interface {
	Source() Unit       // Source of the event
	Name() string       // Name of the event
	NS() gopi.EventNS   // Namespace for the event
	Value() interface{} // Value associated with the event
}
```

Events may be emitted for example by:

  * Your hardware GPIO interface as a pin changes state;
  * A key press, mouse move or button click;
  * A service becoming available on the network;
  * A ticker which fires at a regular interval.

There are many other cases where events could fire. In this chapter, I will describe a tool which handles a ticker, firing at a regular interval.

## The Ticker Unit

Here are the parameters you'll need in order to use the ticker:

| Parameter        | Value                |
| ---------------- | -------------------- |
| Name             | `gopi/timer`         |
| Interface        | `gopi.Timer`         |
| Type             | `gopi.UNIT_TIMER`    |
| Requires         | `gopi.UNIT_BUS`      |
| Import           | `github.com/djthorpe/gopi/v2/unit/timer` |
| Events           | `gopi.TimerEvent`    |
| Compatibility    | Linux, Darwin        |

The interface is defined as follows:

```go
type Timer interface {
	Unit

	NewTicker(time.Duration) gopi.EventId // Create periodic event at interval
	NewTimer(time.Duration) gopi.EventId  // Create one-shot event after interval
	Cancel(gopi.EventId) error            // Cancel events
}
```

You can access the unit instance using the `app.Timer()` convenience method. The relevant `Main` function looks like this:

```go
func Main(app gopi.App, args []string) error {
    app.Timer().NewTicker(time.Second)

	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to exit")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
```

This will fire a `gopi.TimerEvent` once every second. But as nothing has been set up to handle the messages, you may just see some debugging output if you have used the `-debug` flag which indicates the events are not being handled.

In fact, the timer emits the ticker events into a __message bus__, which is unsuprisingly yet another __unit__. You don't need to use the message bus directly, but you can simply use it by defining handlers when setting up your application.

## Setting up event handlers for your application

An event handler is defined as follows:

```go
type EventHandler struct {
	Name string
	Handler EventHandlerFunc
	EventNS EventNS
	Timeout time.Duration
}
```

The fields are:

  * The `Name` of the event;
  * The `Handler` function, which has the signature of `func(context,Context,gopi.App,gopi.Event)`;
  * Optionally, a namespace for events, or `gopi.EVENT_NS_DEFAULT` otherwise;
  * Optionally, a deadline for the event handling, or zero otherwise.

You define an array of event handlers when creating your application. For example,

```go
func TimerHandler(ctx context.Context,app gopi.App, evt gopi.Event) {
	// ...
}

func OtherEventHandler(ctx context.Context,app gopi.App, evt gopi.Event) {
	// ...
}

func Main(app gopi.App, args []string) error {
	// ...
}

func main() {
	Events := []gopi.EventHandler{
		gopi.EventHandler{Name: "gopi.TimerEvent", Handler: TimerHandler},
		gopi.EventHandler{Name: "gopi.Event", Handler: OtherEventHandler},
	}
	app, err := app.NewCommandLineTool(Main, Events, "timer")
	// ...
}
```

The `TimerHandler` method is called whenever a `gopi.TimerEvent` is fired. The `context.Context` parameter can generally be ignored, except for handlers which take a long time to complete, and require a timeout value to cancel. You tool
won't complete until all handlers have returned, so it's possible to create
deadlock.

## Emitting events

When you develop tools you may want to emit your own events. There's an application method `Emit` which does that:

```go

func Main(app gopi.App, args []string) error {
	// ...emit null event...
	app.Emit(gopi.NullEvent)
}
```

The `NullEvent` is simply an event with no information. You can also create your
own events which can be emitted to handlers as long as they adhere to the gopi.Event interface.

Take extra caution when emitting events within handlers. It's easy to create a deadlock situation when you are both handling events of a particular type and also emitting them, resulting in a freeze or panic.

## Conclusion

In the next section you'll see how units like GPIO emit events when the state
changes on a pin (from low to high, or vice-versa, for example). In subsequent chapters (yet to be written!) events can be emitted from input devices, networks or other external stimulous.

The __Message Bus__ is mostly hidden as you respond to events through handlers. But it is implemented through goroutines and channels so using the go language to maximum extent. There are clearly still challenges around syncronization and deadlock to be careful of when developing your own tools and units.


