
# Quickstart Guide

The following sections provide a quickstart guide, with more detailed 
information provided [here](https://godoc.org/github.com/djthorpe/gopi).

## Introduction

This repository contains the Golang interface to hardware and graphics devices
for the Raspberry Pi platform with Linux. The scope of this framework is as 
follows:

  * Developing applications easily, either for use on
    the command line or event-based applications;
  * Enumerating the capabilities of the device and providing various
    information about the device, including the hardware model, serial number
    and so forth;
  * Access to input/output devices, including GPIO, I2C, SPI, Touchscreen,
    Mouse and Keyboard devices;
  * Use of the Graphics Processing Unit (if there is one) including creating
    displays & window surfaces, and being able to move them on the screen;
  * Access to OpenVG and OpenGL API's and rendering text with scalable
    vector fonts;
  * Use of the Multi-Media Abstraction for video, audio & image encoding, 
    decoding and access to the Raspberry Pi camera.

In order to use the library, you'll need to have a working version of Go on 
your Raspberry Pi, which you can [download](https://golang.org/dl/). Then 
retrieve the library on your device, using:

```
  linux% go get github.com/djthorpe/gopi
```

# Running the example programs

There are many examples in the `examples` folder, which can all be installed
on the command line using the following command:

```
  linux% cd "${GOPATH}/src/github.com/djthorpe/gopi"
  linux% examples/build_examples.sh
  linux% helloworld_example
```

Please see the source code for the hello world application 
[here](https://github.com/djthorpe/gopi/tree/master/examples/helloworld).

# What's Next?

Read the remaining documentation on the various functions of `gopi`:

  * To read about getting general information about your device, read [DEVICE](DEVICE.md).
  * To read about input devices, read [INPUT](INPUT.md).
  * To read about opening displays, creating windows and resources, read [EGL](EGL.md).
  * To use OpenVG which provides you with Vector Graphics operations, read [OpenVG](OpenVG.md).
  * To render text as graphics using scalable fonts, read [FONTS](FONTS.md).
  * To use the GPIO peripheral port, read [GPIO](GPIO.md).
  * The application framework can make it easy to write applications. Read [APP](APP.md).

