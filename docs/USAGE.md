
# Organization &amp; Usage

This section describes how the framework is organized and how you can use
it within your own projects. It also describes how you can file feature
requests and issues you are having with the framework.

## Requirements

## Device Independence

The framework is meant to provide some level of device independence, so that
there can be several implementations of the same interfaces. For example,
there are both `rpi.GPIO` and `linux.GPIO` device subsystems which both
implement the `hw.GPIO` interface, but which have different ways to
interact with the GPIO port, different features, etc.

Every device is opened using the `gpio.Open` method with a device-specific
configuration. For the GPIO example,

```
   gpio_rpi, err := gopi.Open(rpi.GPIO{ /* configuration parameters */ },logger)
   if err != nil { /* handle the error */ }
   defer gpio_rpi.Close()

   gpio_linux, err := gopi.Open(linux.GPIO{ /* configuration parameters */ },logger)
   if err != nil { /* handle the error */ }
   defer gpio_linux.Close()
```

Every device driver must conform to the `gopi.Driver` interface which provides
both an `Open` and `Close` method. In order to use the device driver, you also
need to cast the driver to the particular abstract interface (in this case,
a `hw.GPIODriver`. For example,

```
   pins_rpi := gpio_rpi.(hw.GPIODriver).NumberOfPhysicalPins()
   pins_linux := gpio_linux.(hw.GPIODriver).NumberOfPhysicalPins()
```

The abstract device drivers are documented in the following folders:

  * _hw_ The abstract hardware device drivers
  * _kkronos_ The abstract graphics device drivers

The concrete device drivers which implement these are stored in the `device`
folder.

## Organization of the Repository

TODO

## Logging and Debugging

TODO

## Submitting Issues and Feature Requests

TODO

