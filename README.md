# Read me first

| ![GOPI](https://raw.githubusercontent.com/djthorpe/gopi/master/etc/images/gopi-800x388.png) |  Go Language Application Framework |
| :--- | :--- |


[![CircleCI](https://circleci.com/gh/djthorpe/gopi/tree/v3.svg?style=svg)](https://circleci.com/gh/djthorpe/gopi/tree/v3)

This repository contains an application framework for the Go language, which will allow you to develop applications which utilize a number of features of your computer. It's targetted at the Raspberry Pi presently. The following features are intended to be supported:

* The GPIO, I2C and SPI interfaces
* Font loading and rendering in bitmap and vector forms
* Infrared transmission and receiving, for example for remote controls
* Network microservices, announcement and discovery using mDNS and gRPC

It would also be great to support the following features in the future:

* Image and video encoding/decoding, including utilizing hardware

  acceleration

* GPU acceleration for 2D graphics
* 3D graphics
* Audio devices
* Input devices like the mouse, keyboard and touchscreen
* Display and display surfaces, bitmaps and vector graphics
* Connected cameras
* User interface widgets and layout
* Building for Darwin \(Macintosh\) targets

## Requirements

The tested requirements are currently:

* Any Raspberry Pi \(v2, v3, v4, Zero and Zero W have been tested\)
* Raspbian GNU/Linux 9 \(other distributions may work, but not tested\)
* Go 1.13

In order to use the library, you'll need to have a working version of Go on your Raspberry Pi, which you can [download](https://golang.org/dl/). Then retrieve the library on your device, using:

```bash
go get github.com/djthorpe/gopi/v3
```

## Getting Started

In order to get started, build some of the examples in the "cmd" folder. They can be built with the makefile.

* `make rpi` will run the tests and install the examples

Fuller documentation of the examples and developing your own code against this framework will be available in documentation.

## License

```text
Copyright 2016-2020 David Thorpe All Rights Reserved

Redistribution and use in source and binary forms, with or without 
modification, are permitted with some conditions.
```

This repository is released under the BSD License. Please see the file [LICENSE]() for a copy of this license and for a list of the conditions for redistribution and use.

