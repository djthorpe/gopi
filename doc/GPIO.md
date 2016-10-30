
# The GPIO interface

The Raspberry Pi General Purpose interface provides you with a way to interface
with external hardware.

In order to interact with the GPIO device, you'll need to create an object given
the device. For example:

```go
	// Create the GPIO object
    gpio, err := gopi.Open(rpi.GPIO{
		Device: my_device
	}, logger)
	if err != nil { /* handle error */ }
	defer gpio()
	
	/* do things here */
```

You can interact with GPIO "logical pins" by setting their mode to be `INPUT`
or `OUTPUT` and reading the pin level. For example,

```go
	// Return a logical pin from a physical pin and set the mode to output
	pin := gpio.PhysicalPin(40) // GPIO21
	gpio.SetPinMode(pin,gopi.GPIO_OUTPUT)
	
	// Set the pin high
	gpio.WritePin(pin,gopi.GPIO_HIGH)
```

## Understanding logical and physical pins

The physical header on the Raspberry Pi numbers the pins from 1 to 40 (or
1 to 28 for the first version of the Raspberry Pi). Many of these physical
pins "map" to logical pin numbers or names.


