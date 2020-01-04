
# Introduction

__Gopi__ is an application framework for the Go language ("golang"). Ultimately it targets the Raspberry Pi platform running Linux to utilize many of the features of the platform, but it's applicable to any platform where golang runs.

Whilst golang is an excellent programming language to develop with, it needs more work to support multi-media application development. Providing a "framework" with which to develop such applications may speed up their quality and quantity.

The scope of this framework is as follows:

* Allows you to developing applications easily, either for use on the command line or as event-based applications;
* Target several platforms simultaneously, making the best use of the features of that platform;
* Access to input/output devices, including GPIO, I2C, SPI, Touchscreen, Mouse and Keyboard devices (depending on whether that platform supports them);
* Use of the Graphics Processing Unit (if there is one) including creating displays & window surfaces, and being able to move them on the screen;
* Access to vector graphics and font rendering, and basic user interface element layout;
* Remote Procedure Call (RPC) server and client development.

The Raspberry Pi is the first supported hardware platform but it wouldn't be a stretch to provide these facilities on other hardware platforms.

## Requirements

The tested requirements are currently:

  * Any Raspberry Pi (v2, v3, v4, Zero and Zero W have been tested)
  * Raspbian GNU/Linux (Raspian or Buster)
  * Go 1.13

## Building Examples

In order to use the framework, you'll need to have a working version of Go on 
your Raspberry Pi, which you can download from [here](https://golang.org/dl/). Then 
retrieve the library on your device, using:

```sh
bash% git clone https://github.com/djthorpe/gopi
bash% cd gopi
bash% git checkout v2
```

Then, build some of the examples in the "cmd" folder. They can be built with the makefile:

```sh
bash% make test
bash% make install
bash% make clean
```

There are specific versions which can take advantage of Raspberry Pi features:

```sh
bash% make rpi
```

## License

This repository is released under the Apache License. Please see the file
[LICENSE](../LICENSE.md) file in the reposiroty for a copy of this license and 
for a list of the conditions for redistribution and use.

## Contributing

I welcome contributions, feature requests and issues. In order to contribute, please fork
and send pull requests. Feel free to contact me for feature requests and bugs by filing
issues [here](https://github.com/djthorpe/gopi/issues).

## What's Next

The current status is that of a framework is _in development_, with the
features working on some platforms and not others. The following sections are yet to be 
written:

  * Describing Helloworld
  * Event Handling
  * Hardware platforms
  * A section on developing user interfaces, including layout
  * 3D graphics using OpenGL
  * Encoding and decoding video and audio
  * Cameras
  * Developing micro-services
  * Interacting with IR signals from remote controls
  * Developing your own modules and for other platforms
