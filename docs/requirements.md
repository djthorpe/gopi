
The tested requirements are currently:

  * Any Raspberry Pi (v2, v3, v4, Zero and Zero W have been tested)
  * Raspbian GNU/Linux (Raspian or Buster)
  * Go 1.13

In order to use the framework, you'll need to have a working version of Go on 
your Raspberry Pi, which you can [download](https://golang.org/dl/). Then 
retrieve the library on your device, using:

```sh
git clone https://github.com/djthorpe/gopi
```

Then, build some of the examples in the "cmd" folder. They can be built with the makefile.

* `make` will run the tests and install the examples
* `make install` will build and install the examples without testing
* `make clean` will remove intermediate files

