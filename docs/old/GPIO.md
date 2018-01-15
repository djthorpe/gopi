
# The General Purpose Input Output (GPIO) port

The Raspberry Pi General Purpose interface provides you with a way to interface
with external hardware. In order to interact with the GPIO device, you'll need
to create an object given the device. For example:

```go
	// Create the GPIO object
    gpio, err := gopi.Open(rpi.GPIO{
		Device: my_device
	}, logger)
	if err != nil { /* handle error */ }
	defer gpio.Close()
	
	/* do things here */
```

You can interact with GPIO "logical pins" by setting their mode to be `INPUT`
or `OUTPUT` and reading the pin level. For example here is a way to set the
pin GPIO21 to output and set the output to high:

```go
	// Return a logical pin from a physical pin and set the mode to output
	pin := gpio.PhysicalPin(40) // GPIO21
	if pin != gopi.GPIO_PIN_NONE {
		gpio.SetPinMode(pin,gopi.GPIO_OUTPUT)
	}
	
	// Set the pin high
	gpio.WritePin(pin,gopi.GPIO_HIGH)
```

## Understanding logical and physical pins

The physical header on the Raspberry Pi numbers the pins from 1 to 40 (or
1 to 28 for the first version of the Raspberry Pi). Many of these physical
pins "map" to logical pin numbers or names. The mapping could vary from one
device to another, so it's important to check this at runtime. To enumerate
the pins and provide their logical and physical information, use the Pins()
function:

```go
  for _, logical := range gpio.Pins() {
	if physical := gpio.PhysicalPinForPin(pin); physical != 0 {
		fmt.Printf("Logical Pin=%v, Physical Pin=%v\n",logical,physical)
	}
  }
```

## Setting the pin and pull-up internal resistor mode

Pins can be in one of several modes, but it's likely most pins will be in
`INPUT` or `OUTPUT` mode (there are also `ALT0`..`ALT5`). You can set the
pin mode (or query the state of the current mode) using the functions
`SetPinMode` and `GetPinMode`. For example,

```go
	// Return a logical pin from a physical pin and set the mode to output
	pin := gpio.PhysicalPin(40) // GPIO21
	if pin != gopi.GPIO_PIN_NONE {
		gpio.SetPinMode(pin,gopi.GPIO_OUTPUT)
	}
	
	// Print out the current pin mode
	mode := gpio.GetPinMode(pin)
	fmt.Println(pin,"is in mode",mode)
```

The Raspberry Pi also allows you to set the internal pull-up resistor mode to
either `OFF`, `DOWN` or `UP`.

TODO






