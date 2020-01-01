
# Introduction

__Gopi__ is an application framework for the Go language ("golang"). Ultimately it targets the Raspberry Pi platform running Linux to utilize many of the features of the platform, but it's applicable to any platform where golang runs.

The scope of this framework is as follows:

* Allows you to developing applications easily, either for use on the command line or as event-based applications;
* Target several platforms simultaneously, making the best use of the features of that platform;
* Access to input/output devices, including GPIO, I2C, SPI, Touchscreen, Mouse and Keyboard devices (depending on whether that platform supports them);
* Use of the Graphics Processing Unit (if there is one) including creating displays & window surfaces, and being able to move them on the screen;
* Access to vector graphics and font rendering, and basic user interface element layout;
* Remote Procedure Call (RPC) server and client development.

The Raspberry Pi is the first supported hardware platform but it wouldn't be a stretch to provide these facilities on other hardware platforms.

In order to use the framework, you'll need to have a working version of Go which you can download. Then retrieve the framework on your device, using:

```go
  bash% go get github.com/djthorpe/gopi
```

Whilst golang is an excellent programming language to develop with, it needs more work to support multi-media application development. Providing a "framework" with which to develop such applications may speed up their quality and quantity.

The current status is that of a framework is _in development_, with the
features working on some platforms and not others. The framework is split into a number of repositories which are roughly:


| Repository    | Link   | Module |
| ------------- | ------ | ---- |
| gopi-hw       | [`github.com/djthorpe/gopi-hw`](http://github.com/djthorpe/gopi-hw/) | Hardware implementations  |
| gopi-graphics | [`github.com/djthorpe/gopi-graphics`](http://github.com/djthorpe/gopi-graphics/) | Graphics & Fonts |
| gopi-rpc      | [`github.com/djthorpe/gopi-rpc`](http://github.com/djthorpe/gopi-rpc/) | Microservices & Discovery |
| gopi-input    | [`github.com/djthorpe/gopi-input`](http://github.com/djthorpe/gopi-input/) | Input services (Keyboard, Mouse, Touchscreen) |
| gopi-media    | [`github.com/djthorpe/gopi-media`](http://github.com/djthorpe/gopi-media/) | Media services (Encoding, Decoding, Cameras) |


# What's Next?

The following sections are yet to be written:

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
