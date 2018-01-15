
## Introduction & Motivation

This repository implements an application framework for the Go 
programming language ("golang"). Ultimately it targets the 
Raspberry Pi platform running Linux to utilize many of the features 
of the platform, but it's applicable to any platform on which golang runs.

The scope of this framework is as  follows:

  * Allows you to developing applications easily, either for use on
    the command line or as event-based applications;
  * Target several platforms simultaneously, making the best use of the
    features of that platform;
  * Access to input/output devices, including GPIO, I2C, SPI, Touchscreen,
    Mouse and Keyboard devices (depending on whether that platform
    supports them);
  * Use of the Graphics Processing Unit (if there is one) including creating
    displays & window surfaces, and being able to move them on the screen;
  * Access to vector graphics and font rendering, and basic user interface
    element layout.

The [Raspberry Pi](https://www.raspberrypi.org/) is the first supported
hardware platform but it wouldn't be a stretch to provide these facilities on other
hardware platforms.

In order to use the library, you'll need to have a working version of Go which you 
can [download](https://golang.org/dl/). Then retrieve the framework on your 
device, using:

```
  bash% go get github.com/djthorpe/gopi
```

Whilst __golang__ is an excellent programming language to develop with,
it needs more work to support multi-media application development. Providing
a "framework" with which to develop such applications may speed up their
quality and quantity.

# Running the example programs

There are many examples in the `cmd` folder, which can all be installed
on the command line using the following commands:

```
  bash% cd "${GOPATH}/src/github.com/djthorpe/gopi"
  bash% cmd/build_rpi.sh
```

Please see the source code for the hello world application 
[here](https://github.com/djthorpe/gopi/blob/modules/cmd/helloworld/helloworld.go)
to see how a basic command-line application is structured.

# What's Next?

Read the remaining documentation on the various functions of `gopi`:

  * To understand how to use the framework to develop your own applications, see [Helloworld](helloworld.md)
  * Events, tasks and timers are described in [Events](events.md)
  * Information about the hardware platform your applcation us running on is described in [Hardware](hardware.md)

The following sections are yet to be written:

  * A section on developing user interfaces, including layout
  * 3D graphics using OpenGL
  * Encoding and decoding video and audio
  * Cameras
  * Developing micro-services
  * Interacting with IR signals from remote controls
  * Developing your own modules and for other platforms

# Issues, Licensing & Contributing

Contributions and reports of issues are very welcome, using the appropriate __github__ mechanisms.
Otherwise, contact me directly and I hope I can get back to you.

I am using the BSD License which allows you to use the framework in your own software with
attribution. Please see the [LICENSE](https://github.com/djthorpe/gopi/blob/modules/LICENSE.md) for more
information.



 



