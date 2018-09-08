
<table style="border-color: white;"><tr>
  <td width="50%">
    <img src="https://raw.githubusercontent.com/djthorpe/gopi/master/etc/images/gopi-800x388.png" alt="GOPI" style="width:200px">
  </td><td>
    Go Language Application Framework
  </td>
</tr></table>

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
  * Raspian Jessie Lite 4.4 (other distributions may work, but not tested)
  * Go 1.11

In order to use the library, you'll need to have a working version of Go on 
your Raspberry Pi, which you can [download](https://golang.org/dl/). Then 
retrieve the library on your device, using:

```
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

Please see each repository for more example code and information on the _modules_ provided.

# Getting Started

In order to get started, build some of the examples in the "cmd" folder. They
can be built with the makefile.

  * `make` will run the tests and install the examples
  * `make install` will build and install the examples without testing
  * `make clean` will remove intermediate files

Fuller documentation of the examples and developing your own code
against this framework is available in the documentation.

# License

```
Copyright 2016-2018 David Thorpe All Rights Reserved

Redistribution and use in source and binary forms, with or without 
modification, are permitted with some conditions. 
```

This repository is released under the BSD License. Please see the file
[LICENSE](LICENSE.md) for a copy of this license and for a list of the
conditions for redistribution and use.
