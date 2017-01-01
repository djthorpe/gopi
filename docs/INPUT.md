
# Input Devices

In order to obtain input from keyboards, mice, touchscreens and joysticks
there is an abstract interface to input devices which you can use to accept
events such as key presses and positional changes.

There are three main concepts with input devices:

  * The _Input Driver_ provides a way to discover, open, close and watch devices.
  * The _Input Device_ represents a single device such as a mouse, keyboard,
    touchscreen or joystick
  * An _Input Event_ is emitted when an input device changes state. For example,
    a key was pressed or the mouse was moved.
	
The following sections describe how to combine these concepts.

## The Input Driver

You can open a driver in a similar way to other drivers:

```go
	// Create the hw.InputDriver object
	input, err := gopi.Open(linux.InputDriver{ },logger)
	if err != nil { /* handle error */ }
	defer input.Close()
```

To discover a set of devices, simply use the `OpenDevicesByName` method, which
will return devices based on criteria you provide. You can provide a name of
device, a set of device types and the bus on which the device needs to be
connected. These arguments are described in more detail below.

The function will return an array of new devices opened and any error code if the
operation was unsuccessful. If you call this method more than once, it will only
open devices if they are not yet opened.

Looking at the arguments to the method, you can set the `name` parameter as
empty or refer to a device by it's full name, alias or physical connection
name. If empty, all devices discovered will be considered.

The `flags` argument can either be `hw.INPUT_TYPE_ANY` to match all types of
device, or you can provide an OR'ed set of device types. For example, to open
any mouse and keyboard:

```go
  devices, err := input.(hw.InputDriver).OpenDevicesByName("",hw.INPUT_TYPE_MOUSE | hw.INPUT_TYPE_KEYBOARD,hw.INPUT_BUS_ANY)
  if err != nil {
	/* handle the error */
  }
  if len(devices) == 0 {
    /* no devices found */
  }
```

The `bus` argument can be used to open devices on a specific bus, or use
`hw.INPUT_BUS_ANY` otherwise.

To close devices (which means they are no longer polled for events) use the
`CloseDevice` method:

```go
  err := input.(hw.InputDriver).CloseDevice(mouse)
  if err != nil { /* handle the error */ }
```

Finally, for all opened devices, you will want to watch for events emitted by
the devices. The `Watch` method uses a `delta` argument which determines how
long you want to watch a device for, and a `callback` method which is called
for each event emitted. If no events are emitted within the specified time, then
the Watch method returns `nil`:

```go
  err := input.(hw.InputDriver).Watch(time.Second * 5,func (event hw.InputEvent,device hw.InputDevice) {
	fmt.Println("DEVICE",device)
	fmt.Println("EVENT",event)
  })
  if err != nil { /* handle the error */ }
```

Practically, you will want to continue watching for events until your application
ends, and you'll want to do it in the background. See below for information on
how to implement this pattern.

## Input Events

The `hw.InputEvent` structure provides information on the event which has been
emitted:

  TODO

There are a number of different kinds of events which are emitted and which
you can handle:

  * `INPUT_EVENT_KEYPRESS` - A key or mouse button was pressed down
  * `INPUT_EVENT_KEYRELEASE` - A key or mouse button was released
  * `INPUT_EVENT_KEYREPEAT` - A key continues to be pressed and initiates a repeat
  * `INPUT_EVENT_ABSPOSITION` - The absolute position for the device was changed (for example, a touchscreen was pressed)
  * `INPUT_EVENT_RELPOSITION` - The relative position for the device was changed (for example, a mouse was moved)

For key presses, the key code is provided as part of the event. For
absolute and relative position changes, the `Position` parameter is set, and
for relative position changes, the `Relative` parameter is also set.

## Input Devices

TODO



