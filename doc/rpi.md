
# rpi

The `rpi` package contains the interface to the Raspberry Pi. As well as
providing you with information about the Raspberry Pi, you can create 
GPIO and I2C interfaces.

## Creating a Raspberry Pi object

In order to control the Raspberry Pi, you'll need to create a device object
as follows:

```go

import "gopi/rpi"

func main() {
  device, err := rpi.New()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error: ", err)
    os.Exit(-1)
  }
  defer device.Close()
}
```
You can then query the Raspberry Pi for various information, for example:

```go
  warranty, err := device.WarrantyBit()
  product, err := device.Product() /* rpi.Product */
  processor, err := device.Processor() /* rpi.Processor */
  product_name, err := device.ProductName()
  processor_name, err := device.ProcessorName()
  peripheral_base, err := device.PeripheralBase()
```

Here are the current possible values for Product and Processor:

```go
	RPI_MODEL_UNKNOWN
	RPI_MODEL_A
	RPI_MODEL_B
	RPI_MODEL_A_PLUS
	RPI_MODEL_B_PLUS
	RPI_MODEL_B_PI_2
	RPI_MODEL_B_PI_3
	RPI_MODEL_ALPHA
	RPI_MODEL_COMPUTE_MODULE
	RPI_MODEL_ZERO
```

```go
	RPI_PROCESSOR_UNKNOWN
	RPI_PROCESSOR_BCM2835
	RPI_PROCESSOR_BCM2836
	RPI_PROCESSOR_BCM2837
```

## Retreiving information from VideoCore

You can retrieve information from the VideoCore using the `VCGenCmd` method
and some utility functions:

```go

import "gopi/rpi"

func main() {
  device, err := rpi.New()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error: ", err)
    os.Exit(-1)
  }
  defer device.Close()
  
  response, err := device.VCGenCmd("otp_dump")
  fmt.Println(response)
  
}
```

A list of all defined command names can be retrived using the `GetCommands`
method, which returns an array of command strings. Additional utility functions 
are defined as convenience methods:

  * `commands, err := device.GetCommands()`
  * `temperature, err := device.GetCoreTemperatureCelcius()`
  * `frequencies, err := device.GetClockFrequencyHertz()`
  * `volts, err := device.GetVolts()`
  * `codecs, err := device.GetCodecs()`
  * `memory_size, err := device.GetMemoryMegabytes()`
  * `memory, err := device.GetOTP()`
  * `serial_number, err := device.GetSerial()`
  * `revision, err := device.GetRevision()`
  
