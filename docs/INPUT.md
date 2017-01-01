
# Input Devices

In order to obtain input from keyboards, mice, touchscreens and joysticks
there is an abstract interface to input devices which you can use to accept
events such as key presses and positional changes. The following sections
describe how to combine these concepts.

## Abstract Interface

There are three main concepts with input devices:

  * The **Input Driver** provides a way to discover, open, close and watch devices.
  * The **Input Device** represents a single device such as a mouse, keyboard,
    touchscreen or joystick
  * An **Input Event** is emitted when an input device changes state. For example,
    a key was pressed or the mouse was moved.

There are other enumerations, interfaces and structs which support these concepts:	

| **Import** | `github.com/djthorpe/gopi/hw` |
| -- | -- | -- |
| **Interface** | `hw.InputDriver` | gopi.Driver, implements driver for all devices |
| **Interface** | `hw.InputDevice` | gopi.Driver, implements a single device |
| **Enum**   | `hw.InputDeviceType` | Type of an input device |
| **Enum**   | `hw.InputDeviceBus` | How the input device is connected |
| **Struct** | `hw.InputEvent` | An event emitted by a device |
| **Enum**   | `hw.InputEventType` | Type of event emitted |
| **Enum**   | `hw.InputKeyCode` | The key pressed or mouse button activated |
| **Enum**   | `hw.InputKeyState` | Keyboard state (Caps Lock, Num Lock, Shift, etc) |
| **Function** | `hw.InputEventCallback` | func(event *hw.InputEvent, device hw.InputDevice) |

## Concrete Implementation

The concrete implementation of the driver and devices are currently only for linux:

| **Import** | `github.com/djthorpe/gopi/device/linux` |
| -- | -- | -- |
| **Struct** | `linux.InputDriver` | Concrete Linux input driver configuration |
| **Struct** | `linux.InputDevice` | Concrete Linux input device |

## The Input Driver

The `linux.InputDriver` configuration supports a single configuration parameter:

| **Struct** | `linux.InputDriver` |
| -- | -- | -- |
| **Bool** | Exclusive | Whether to open devices for exclusive access |

The input driver and input device implements the `gopi.Driver` interface so
can be opened in the usual way:

```go
input, err := gopi.Open(linux.InputDriver{ Exclusive: true },logger)
if err != nil { /* handle error */ }
defer input.Close()
```

The `hw.InputDriver` interface should implement the following methods:

| **Interface** | `hw.InputDriver` |
| -- | -- | -- |
| **Method** | `Close() error` | Release all devices and close |
| **Method** | `OpenDevicesByName(name string, flags hw.InputDeviceType, bus hw.InputDeviceBus) ([]hw.InputDevice, error)` | Open devices |
| **Method** | `CloseDevice(device hw.InputDevice) error` | Close a device |
| **Method** | `Watch(delta time.Duration,callback hw.InputEventCallback) error` | Watch for events and callback on emitted event |

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
device, or you can provide an OR'ed set of device types.

| **Enum** | `hw.InputDeviceType` |
| -- | -- |
| `hw.INPUT_TYPE_NONE` | None or unknown device type |
| `hw.INPUT_TYPE_KEYBOARD` | Keyboard |
| `hw.INPUT_TYPE_MOUSE` | Mouse |
| `hw.INPUT_TYPE_TOUCHSCREEN` | Multi-touch input device |
| `hw.INPUT_TYPE_ANY` | Matches any device when calling `OpenDevicesByName` |

For example, to open any mouse and keyboard:

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

| **Enum** | `hw.InputDeviceBus` |
| -- | -- |
| `hw.INPUT_BUS_NONE` | Unknown device bus |
| `hw.INPUT_BUS_USB` | USB Bus |
| `hw.INPUT_BUS_BLUETOOTH` | Bluetooth Bus |
| `hw.INPUT_BUS_ANY` | Matches any bus when calling `OpenDevicesByName` |

To close devices use the `CloseDevice` method:

```go
  err := input.(hw.InputDriver).CloseDevice(mouse)
  if err != nil { /* handle the error */ }
```

By closing a device, it is removed from the list of opened devices,
the device no longer emits events and any exclusivity for access is released.

## Input Events

For all opened devices, events are emitted and can be consumed by calling
the `Watch` method. This method uses a `delta` argument which determines how
long you want to watch a device for, and a `callback` method which is called
for each event emitted. If no events are emitted within the specified time, then
the Watch method returns `nil`. It will return with an error before this if
some unexpected condition occurs:

```go
  err := input.(hw.InputDriver).Watch(time.Second * 5,func (event hw.InputEvent,device hw.InputDevice) {
	fmt.Println("DEVICE=",device)
	fmt.Println("EVENT=",event)
  })
  if err != nil { /* handle the error */ }
```

Practically, you will want to continue watching for events until your application
ends, and you'll want to do it in the background. See below for information on
how to implement this pattern.

The `hw.InputEvent` structure provides information on the event which has been
emitted:

| **Struct** | `hw.InputEvent` |
| -- | -- | -- |
| **time.Duration** | `Timestamp` | The timestamp for an event, to nanosecond resolution |
| **hw.InputDeviceType** | `DeviceType` | The type of device emitting the event |
| **hw.EventType** | `EventType` | The type of event being emitted |
| **hw.InputKeyCode** | `Keycode` | For press, release and repeat events, the key or mouse button |
| **uint32** | `Scancode` | The keyboard scancode for press and release events |
| **khronos.EGLPoint** | `Position` | For touchscreen and mouse device types, the absolute position being tracked |
| **khronos.EGLPoint** | `Relative` | For mouse device types, the relative position compared to last position |
| **uint** | `Slot` | For multi-touch events, the slot number |

There are a number of different types of events. Fields of `hw.InputEvent` are
populated differently depending on the type of event.

| **Enum** | `hw.InputEventType` |
| -- | -- |
| `hw.INPUT_EVENT_KEYPRESS` | Mouse, touchscreen or keyboard key press |
| `hw.INPUT_EVENT_KEYRELEASE` | Mouse, touchscreen or keyboard key release |
| `hw.INPUT_EVENT_KEYREPEAT` | Keyboard key being held down |
| `hw.INPUT_EVENT_ABSPOSITION` | Mouse or touch screen absolute position |
| `hw.INPUT_EVENT_RELPOSITION` | Mouse or touch screen relative position change |
| `hw.INPUT_EVENT_TOUCHPRESS` | Touchscreen press |
| `hw.INPUT_EVENT_TOUCHRELEASE` | Touchscreen press |
| `hw.INPUT_EVENT_TOUCHPOSITION` | Touchscreen position change |

For **keyboard** devices, the `Keycode` and `Scancode` fields will be set, where the
key code determines the pressed key. Scancode is usually a device-specific
code translated into the keycode by the keyboard hardware.

For **mouse* devices, relative positions are reported, but the absolute position
is also set synthetically and can be changed at any time by calling the
`SetPosition` method of `hw.InputDevice`. Mouse devices also report on button
presses.

For **touchscreen** devices, these are generally _multi-touch_. For example, you
can activate two different points on the touch device simultaneously with fingers
or styluses. Touchscreen devices will therefore not only report on absolute
position changes and key press events, but also report individually using the
`hw.INPUT_EVENT_TOUCHPRESS`, `hw.INPUT_EVENT_TOUCHRELEASE` and `hw.INPUT_EVENT_TOUCHPRESS`
for each of the simultaneous touches, or _slots_.

## Input Devices

TODO



