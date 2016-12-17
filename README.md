# Introduction

This repository contains the Golang interface to hardware and graphics devices,
currently targetting the Raspberry Pi but could theoretically target other
platforms too.

The scope of this library is as follows:

  * A framework for developing applications easily, either for use on
    the command line or event-based applications
  * Enumerating the capabilities of the device and providing various
    information about the device, including the hardware model, serial number
    and so forth
  * Access to input/output devices, including GPIO, I2C, SPI, Touchscreen,
    Mouse and Keyboard devices
  * Use of the Graphics Processing Unit (if there is one) including creating
    displays & window surfaces, and being able to move them on the screen
  * Access to OpenVG and OpenGL API's

In order to use the library, you'll need to have a working version of Go on 
your Raspberry Pi, which you can [download](https://golang.org/dl/). Then 
retrieve the library on your device, using:

```
go get github.com/djthorpe/gopi
```

The following sections provide a quickstart guide, with more detailed 
information provided [here](https://godoc.org/github.com/djthorpe/gopi).

# Quickstart Guide

The `gopi` package abstracts the Raspberry Pi device, peripherals and GPU. It
also includes a logger interface which is used throughout. Here's a short 
example which creates a logger and then prints out some information about 
the Raspberry Pi:

```go

package main

import (
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi/util"
	"github.com/djthorpe/gopi/device/rpi"
)

func main() {
	
	// Create the logger
	log, err := util.Logger(util.StderrLogger{ })
	if err != nil {
		panic("Can't open logging interface")
	}
	defer log.Close()

	// Set logging level
	log.SetLevel(util.LOG_ANY)

	// Open the Raspberry Pi device
	device, err := gopi.Open(rpi.Device{ },log)
	if err != nil {
		log.Fatal("%v",err)
		return
	}
	defer device.Close()

	log.Info("Device=%v",device)
}

```
The `gopi.Open` method takes two arguments, one of which is your "concrete"
Raspberry Pi device, and the second of which is a reference to your logger.
If you leave the second argument is `nil` then a logging device is created
for you, but it doesn't print anything out.

You can get some information on your Raspberry Pi using the following method:

```go
  log.Info("Info=%v",device)
```

This displays model number, serial number, processor, PCB revision and
so forth.

# Running the example programs

There are lots of example programs in the `examples` folder. You can run any
of these using the `go run` command or compile and install them using `go install`.
The examples demonstrate various features of the `gopi` package. For example:

  * `helloworld_example.go` Simply displays the canonical simplest program in Go
  * `log_example.go` Demonstates how to write your own logger class. In this case
      it logs to a file rather than to the console.
  * `windowing_example.go` Demonstates how to draw onto the screen using the
      GPU.
  * `gpio_example.go` Demonstrates how to use the GPIO interface on the Raspberry Pi.
  * `led_example.go` Demonstrates how to use the GPIO interface on the Raspberry Pi
	to blink an LED.

# What's Next?

Read the remaining documentation on the various functions of `gopi`:

  * To read about getting general information about your device, read [DEVICE](doc/DEVICE.md).
  * To read about input devices, read [INPUT](doc/INPUT.md).
  * To read about opening displays, creating windows and resources, read [EGL](doc/EGL.md).
  * To use OpenVG which provides you with Vector Graphics operations, read [OpenVG](doc/OpenVG.md).
  * To render text as graphics using scalable fonts, read [Fonts](doc/FONTS.md).
  * To use the GPIO peripheral port, read [GPIO](doc/GPIO.md).
  * The application framework can make it easy to write applications. Read [APP](doc/APP.md).

# License

This repository is released under the BSD License. Please see the file
[LICENSE](LICENSE.md) for a copy of this license.

