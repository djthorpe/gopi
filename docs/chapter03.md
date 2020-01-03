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

There are many other cases where events could fire. In this chapter, I will
describe a tool which handles a ticker, firing at a regular interval.

## The Ticker unit

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

This will fire a `gopi.Event` once every second. But as nothing has been set up to handle the messages, you may just see some debugging output if you have used the `-debug` flag which indicates the events are not being handled.

In fact, the timer emits the ticker events into the __message bus__, which is unsuprisingly yet another __unit__.

## The Message Bus unit

The message bus allows any other unit to:

  * Emit events
  * Register functions that can handle events
  * Register a default function which can handle any events not handled otherwise.

Here are the parameters you'll need in order to access the message bus:

| Parameter        | Value                |
| ---------------- | -------------------- |
| Name             | `gopi/bus`           |
| Interface        | `gopi.Bus`           |
| Type             | `gopi.UNIT_BUS`      |
| Import           | `github.com/djthorpe/gopi/v2/unit/bus` |
| Compatibility    | Linux, Darwin        |

The interface is defined as follows:

```go
type Bus interface {
	Unit

	Emit(gopi.Event) // Emit an event on the bus
    NewHandler(string, gopi.EventHandler) error 	// Register an event handler
    DefaultHandler(EventNS, gopi.EventHandler) error // Register a default handler
}
```

You can access the unit instance using the `app.Bus()` convenience method. The relevant `Main` function then looks like this:

```go
func HandleTicker(_ context.Context,evt gopi.Event) {
    fmt.Println("Event=",evt.(gopi.TimerEvent))
}

func Main(app gopi.App, args []string) error {
    app.Bus().NewHandler("gopi.TimerEvent",HandleTicker)
    app.Timer().NewTicker(time.Second)

	// Wait for CTRL+C
	fmt.Println("Press CTRL+C to exit")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}
```

The `HandleTicker` method is called whenever an appropriate event is fired, in this case the `gopi.TimerEvent`. The `context.Context` parameter can generally be ignored, except for handlers which take a long time to complete.

