
<table style="border-color: white;"><tr>
  <td width="50%">
    <img src="https://raw.githubusercontent.com/djthorpe/gopi/master/etc/images/gopi-800x388.png" alt="GOPI" style="width:200px">
  </td><td>
    Go Language Application Framework
  </td>
</tr></table>

[![CircleCI](https://circleci.com/gh/djthorpe/gopi/tree/master.svg?style=svg)](https://circleci.com/gh/djthorpe/gopi/tree/master)

This repository contains an application framework for the Go language, which
will allow you to develop applications which utilize a number of features
of your computer. It's targetted at the Raspberry Pi presently. The following
features are intended to be supported:

  * The GPIO, I2C and SPI interfaces
  * Display and display surfaces, bitmaps and vector graphics
  * GPU acceleration for 2D graphics
  * Font loading and rendering in bitmap and vector forms
  * Input devices like the mouse, keyboard and touchscreen
  * Infrared transmission and receiving, for example for remote controls
  * Network microservices, announcement and discovery

It would also be great to support the following features in the future:

  * Image and video encoding/decoding, including utilizing hardware
    acceleration
  * Connected cameras
  * 3D graphics
  * Audio devices
  * User interface widgets and layout
  * Building for Darwin (Macintosh) targets

# Requirements

The tested requirements are currently:

  * Any Raspberry Pi (v2, v3, Zero and Zero W have been tested)
  * Raspbian GNU/Linux 9 (other distributions may work, but not tested)
  * Go 1.12

In order to use the library, you'll need to have a working version of Go on 
your Raspberry Pi, which you can [download](https://golang.org/dl/). Then 
retrieve the library on your device, using:

```sh
go get github.com/djthorpe/gopi
```

The `gopi` repository is mostly just a set of interfaces and some utility packages.
Other repositories provide implementation:

| Repository    | Link   | Module |
| ------------- | ------ | ---- |
| gopi-hw       | [`http://github.com/djthorpe/gopi-hw/`](http://github.com/djthorpe/gopi-hw/) | Hardware implementations  |
| gopi-graphics | [`http://github.com/djthorpe/gopi-graphics/`](http://github.com/djthorpe/gopi-graphics/) | Graphics & Fonts |
| gopi-rpc      | [`http://github.com/djthorpe/gopi-rpc/`](http://github.com/djthorpe/gopi-rpc/) | Microservices & Discovery |
| gopi-input    | [`http://github.com/djthorpe/gopi-input/`](http://github.com/djthorpe/gopi-input/) | Input services (Keyboard, Mouse, Touchscreen) |

Please see each repository for more example code and information on the modules provided.

# Getting Started

In order to get started, build some of the examples in the "cmd" folder. They
can be built with the makefile.

  * `make` will run the tests and install the examples
  * `make install` will build and install the examples without testing
  * `make clean` will remove intermediate files

There are two examples in this repository you can examine in order so you can
develop your own applications:

  * `helloworld` demonstrates the most canonical code, taking in command-line
    arguments, outputting a message and waiting for user input to end the
    program;
  * `timers` demonstrates the use of the timer module, either outputting a
    single message, or one on a repeating basis.

Fuller documentation of the examples and developing your own code against this 
framework is available in the documentation.

# Modules

This repository contains two modules:

| Module | Import | Type | Name |
| -------- | ------ | ---- | ---- |
| Logger | `github.com/djthorpe/gopi/sys/logger` | `gopi.MODULE_TYPE_LOGGER` | `sys/logger` |
| Timer | `github.com/djthorpe/gopi/sys/timer` | `gopi.MODULE_TYPE_TIMER` | `sys/timer` |

## Logger

The logger module provides very basic logging functionality. Here is the interface for any
logging module:


```go
type Logger interface {
  Driver

  // Output logging messages
  Fatal(format string, v ...interface{}) error
  Error(format string, v ...interface{}) error
  Warn(format string, v ...interface{})
  Info(format string, v ...interface{})
  Debug(format string, v ...interface{})
  Debug2(format string, v ...interface{})

  // Return IsDebug flag
  IsDebug() bool
}
```

A logger module is **required** for every application which uses this framework, so 
include the module in your main package:

```go
package main

import (
  // Frameworks
  "github.com/djthorpe/gopi"

  // Modules
  _ "github.com/djthorpe/gopi/sys/logger"
)
```

The logger is then available as `app.Logger` within your application, and is also passed
to every gopi module with the `Open` method. The standard logger includes some command-line
flags in case you want to log to a file, rather than to `stderr`:

```go
$ helloworld -help
Usage of helloworld:
  -debug
    	Set debugging mode
  -log.append
    	When writing log to file, append output to end of file
  -log.file string
    	File for logging (default: log to stderr)
  -verbose
    	Verbose logging
```

Logging occurs depending on the combination of the `debug` and `verbose` flags, according to
the following rules:

| Debug   | Verbose | Levels logged         |
| ------- | ------- | --------------------- |
| `false` | `false` | Fatal, Error and Warn |
| `false` | `true`  | Fatal, Error, Warn and Info |
| `true`  | `false` | Fatal, Error, Warn, Info and Debug |
| `true`  | `true`  | Fatal, Error, Warn, Info, Debug and Debug2 |


## Timer

The timer module emits `gopi.Event` objects once, at regular intervals,
or at intervals according to a backoff rule. The timer interface is as follows:

```go
type Timer interface {
  Driver
  Publisher

  // Schedule a timeout (one shot)
  NewTimeout(duration time.Duration, userInfo interface{})

  // Schedule an interval, which can fire immediately
  NewInterval(duration time.Duration, userInfo interface{}, immediately bool)

  // Schedule a backoff timer with maximum backoff
  NewBackoff(duration time.Duration, max_duration time.Duration, userInfo interface{})
}
```

You can subscribe to emitted events which are as follows:

```go
type TimerEvent interface {
  Event

  // Provide the timestamp for the event
  Timestamp() time.Time

  // The user info for the event
  UserInfo() interface{}

  // Cancel the timer which fired this event
  Cancel()
}
```

# License

```
Copyright 2016-2018 David Thorpe All Rights Reserved

Redistribution and use in source and binary forms, with or without 
modification, are permitted with some conditions. 
```

This repository is released under the BSD License. Please see the file
[LICENSE](LICENSE.md) for a copy of this license and for a list of the
conditions for redistribution and use.
