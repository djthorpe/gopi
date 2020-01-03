# Event Handling

Fundamentally the `gopi` framework implements tools by reacting to events. The `gopi.Event` interface defines a basic event to be handled:

```go
type Event interface {
	Source() Unit       // Source of the event
	Name() string       // Name of the event
	NS() EventNS        // Namespace for the event
	Value() interface{} // Value associated with the event
}
```

Events may be emitted for example by:

  * You hardware GPIO interface as a pin changes state;
  * A key press, mouse move or button click;
  * A service becoming available on the network;
  * A ticker which fires at a regular interval.

There are many other cases where events could fire. In this chapter, I will
desxcribe a tool which handles a ticker, firing at a regular interval.

## The Ticker unit

TODO
