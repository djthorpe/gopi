# Introduction

This repository contains Golang interface to the Raspberry Pi hardware, which
is one of the Broadcom ARM processors, which should eventually include
the peripheral devices and the VideoCore GPU. In order to use it, you'll need
to have a working version of Go on your Raspberry Pi, which you
can [download](https://golang.org/dl/). Then in order to retrieve the source 
code on your device, use:

```
go get github.com/djthorpe/gopi
```

# Quickstart Guide

The `gopi` package abstracts the Raspberry Pi device, peripherals and GPU. It
also includes a logger interface which is used throughout in order to provide
informational information. Here's a short example which creates a logger
and then prints out some information about the Raspberry Pi:

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

	// DO STUFF HERE
	log.Info("Device=%v",device)
}

```

The `gopi.Open` method takes two arguments, one of which is your "concrete"
Raspberry Pi device, and the second of which is a reference to your logger.
If you leave the second argument is `nil` then a logging device is created
for you, but it doesn't print anything out.

You can get some information on your Raspberry Pi using the following method:

```go
  log.Info("Info=%v",device.GetInfo())
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

# What's Next?

Read the remaining documentation on the various functions of `gopi`:

  * To read about getting general information about your device, read [DEVICE](doc/DEVICE.md).
  * To read about input devices, read [INPUT](doc/INPUT.md).
  * To read about opening displays, creating windows and resources, read [EGL](doc/EGL.md).
  * To use OpenVG which provides you with Vector Graphics operations, read [OpenVG](doc/OpenVG.md).
  * To use the GPIO peripheral port, read [GPIO](doc/GPIO.md).

# License

This repository is released under the BSD License. Please see the file
[LICENSE](LICENSE.md) for a copy of this license.

# Appendices

TODO
